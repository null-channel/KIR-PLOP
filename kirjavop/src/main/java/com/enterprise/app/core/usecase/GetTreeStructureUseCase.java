package com.enterprise.app.core.usecase;

import com.enterprise.app.core.domain.model.TreeNode;
import com.enterprise.app.core.port.out.KubernetesPort;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
@RequiredArgsConstructor
public class GetTreeStructureUseCase 
        implements UseCase<GetTreeStructureUseCase.InputValues, GetTreeStructureUseCase.OutputValues> {

    private final KubernetesPort kubernetesPort;

    @Override
    public OutputValues execute(InputValues input) {
        List<TreeNode> pods = kubernetesPort.getTreePods(input.namespace);
        TreeNode root = kubernetesPort.constructBTreeFromPods(pods);
        return new OutputValues(root);
    }

    public record InputValues(String namespace) implements UseCase.InputValues {
    }

    public record OutputValues(TreeNode root) implements UseCase.OutputValues {
    }
} 