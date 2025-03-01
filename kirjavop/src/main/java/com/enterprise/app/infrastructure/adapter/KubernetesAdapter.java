package com.enterprise.app.infrastructure.adapter;

import com.enterprise.app.core.domain.model.TreeNode;
import com.enterprise.app.core.port.out.KubernetesPort;
import io.kubernetes.client.openapi.ApiClient;
import io.kubernetes.client.openapi.ApiException;
import io.kubernetes.client.openapi.apis.CoreV1Api;
import io.kubernetes.client.openapi.models.V1Pod;
import io.kubernetes.client.openapi.models.V1PodList;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;

import java.util.ArrayList;
import java.util.Comparator;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;
import java.util.HashMap;

@Slf4j
@Component
@RequiredArgsConstructor
public class KubernetesAdapter implements KubernetesPort {

    private final ApiClient kubernetesClient;

    @Override
    public List<TreeNode> getTreePods(String namespace) {
        try {
            CoreV1Api api = new CoreV1Api(kubernetesClient);
            V1PodList podList = api.listNamespacedPod(
                namespace,
                null,
                null,
                null,
                null,
                "app=tree",  // Label selector for tree pods
                null,
                null,
                null,
                null,
                null
            );

            return podList.getItems().stream()
                .map(this::convertPodToTreeNode)
                .collect(Collectors.toList());
        } catch (ApiException e) {
            log.error("Failed to get pods from Kubernetes", e);
            throw new RuntimeException("Failed to get pods from Kubernetes", e);
        }
    }

    @Override
    public TreeNode constructBTreeFromPods(List<TreeNode> pods) {
        if (pods == null || pods.isEmpty()) {
            return null;
        }

        // Sort pods by left value to ensure proper tree construction
        pods.sort(Comparator.comparingInt(TreeNode::getLeftValue));

        // Find the root node (smallest left value)
        TreeNode root = pods.get(0);
        List<TreeNode> remainingNodes = new ArrayList<>(pods.subList(1, pods.size()));

        // Recursively build the tree
        buildTree(root, remainingNodes);

        return root;
    }

    private void buildTree(TreeNode parent, List<TreeNode> nodes) {
        List<TreeNode> children = nodes.stream()
            .filter(node -> isChild(parent, node))
            .collect(Collectors.toList());

        for (TreeNode child : children) {
            parent.addChild(child);
            nodes.remove(child);
            buildTree(child, nodes);
        }
    }

    private boolean isChild(TreeNode parent, TreeNode potentialChild) {
        return potentialChild.getLeftValue() > parent.getLeftValue() 
            && potentialChild.getRightValue() < parent.getRightValue();
    }

    private TreeNode convertPodToTreeNode(V1Pod pod) {
        Map<String, String> labels = pod.getMetadata().getLabels();
        return new TreeNode(
            pod.getMetadata().getName(),
            Integer.parseInt(labels.getOrDefault("l", "0")),
            Integer.parseInt(labels.getOrDefault("r", "0"))
        );
    }

    public TreeNode buildOperatorTreeStructure(String namespace) {
        List<TreeNode> nodes = getOperatorManagedPods(namespace);
        Map<Integer, TreeNode> nodeMap = createNodeMap(nodes);
        connectNodes(nodes, nodeMap);
        return findRootNode(nodes);
    }

    private List<TreeNode> getOperatorManagedPods(String namespace) {
        try {
            CoreV1Api api = new CoreV1Api(kubernetesClient);
            V1PodList podList = api.listNamespacedPod(
                namespace,
                null,
                null,
                null,
                null,
                null,
                null,
                null,
                null,
                null,
                null
            );

            return podList.getItems().stream()
                .filter(pod -> pod.getMetadata().getOwnerReferences() != null
                    && pod.getMetadata().getOwnerReferences().stream()
                        .anyMatch(ref -> "Tree".equals(ref.getKind())))
                .map(this::convertOperatorPodToTreeNode)
                .collect(Collectors.toList());
        } catch (ApiException e) {
            log.error("Failed to get pods from Kubernetes", e);
            throw new RuntimeException("Failed to get pods from Kubernetes", e);
        }
    }

    private TreeNode convertOperatorPodToTreeNode(V1Pod pod) {
        String name = pod.getMetadata().getName();
        Map<String, String> labels = pod.getMetadata().getLabels();
        int leftChild = labels != null ? Integer.parseInt(labels.getOrDefault("l-child", "-1")) : -1;
        int rightChild = labels != null ? Integer.parseInt(labels.getOrDefault("r-child", "-1")) : -1;
        return new TreeNode(name, leftChild, rightChild);
    }

    private Map<Integer, TreeNode> createNodeMap(List<TreeNode> nodes) {
        Map<Integer, TreeNode> nodeMap = new HashMap<>();
        for (int i = 0; i < nodes.size(); i++) {
            nodeMap.put(i + 1, nodes.get(i));
        }
        return nodeMap;
    }

    private void connectNodes(List<TreeNode> nodes, Map<Integer, TreeNode> nodeMap) {
        for (TreeNode node : nodes) {
            if (node.getLeftValue() >= 0) {
                TreeNode leftChild = nodeMap.get(node.getLeftValue());
                if (leftChild != null) {
                    node.addChild(leftChild);
                }
            }
            if (node.getRightValue() >= 0) {
                TreeNode rightChild = nodeMap.get(node.getRightValue());
                if (rightChild != null) {
                    node.addChild(rightChild);
                }
            }
        }
    }

    private TreeNode findRootNode(List<TreeNode> nodes) {
        return nodes.stream()
            .filter(node -> nodes.stream()
                .noneMatch(parent -> 
                    parent.getLeftValue() == nodes.indexOf(node) + 1
                    || parent.getRightValue() == nodes.indexOf(node) + 1))
            .findFirst()
            .orElse(null);
    }
} 