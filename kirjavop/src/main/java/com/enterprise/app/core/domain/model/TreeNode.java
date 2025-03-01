package com.enterprise.app.core.domain.model;

import lombok.Builder;
import lombok.Getter;
import lombok.Setter;

import java.util.ArrayList;
import java.util.List;

@Getter
@Setter
@Builder
public class TreeNode {
    private String podName;
    private int leftValue;
    private int rightValue;
    private List<TreeNode> children;

    public TreeNode(String podName, int leftValue, int rightValue) {
        this.podName = podName;
        this.leftValue = leftValue;
        this.rightValue = rightValue;
        this.children = new ArrayList<>();
    }

    public TreeNode(String podName, int leftValue, int rightValue, List<TreeNode> children) {
        this.podName = podName;
        this.leftValue = leftValue;
        this.rightValue = rightValue;
        this.children = children;
    }

    public void addChild(TreeNode child) {
        if (children == null) {
            children = new ArrayList<>();
        }
        children.add(child);
    }

    public boolean isLeaf() {
        return children == null || children.isEmpty();
    }
} 