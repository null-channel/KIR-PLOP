package com.enterprise.app.infrastructure.web;

import com.enterprise.app.core.usecase.GetTreeStructureUseCase;
import com.enterprise.app.infrastructure.adapter.KubernetesAdapter;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.Parameter;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.client.RestTemplate;
import com.enterprise.app.core.domain.model.TreeNode;

@RestController
@RequestMapping("/v1/trees")
@RequiredArgsConstructor
@Tag(name = "Tree Operations", description = "API endpoints for tree operations")
public class TreeController {

    private final GetTreeStructureUseCase getTreeStructureUseCase;
    private final KubernetesAdapter kubernetesAdapter;
    private final RestTemplate restTemplate;

    @GetMapping("/{namespace}")
    @Operation(summary = "Get tree structure from Kubernetes pods",
              description = "Retrieves pods with tree labels and constructs a B-tree based on 'l' and 'r' tags")
    public ResponseEntity<TreeNode> getTreeStructure(
            @Parameter(description = "Kubernetes namespace to search for tree pods")
            @PathVariable String namespace) {
        
        GetTreeStructureUseCase.InputValues input = new GetTreeStructureUseCase.InputValues(namespace);
        GetTreeStructureUseCase.OutputValues output = getTreeStructureUseCase.execute(input);
        
        return ResponseEntity.ok(output.root());
    }

    @GetMapping("/operator/{namespace}")
    @Operation(summary = "Get tree structure from operator-managed pods",
              description = "Retrieves pods managed by the Tree operator and constructs their tree structure")
    public ResponseEntity<TreeNode> getOperatorTreeStructure(
            @Parameter(description = "Kubernetes namespace to search for operator-managed tree pods")
            @PathVariable String namespace) {
        
        TreeNode root = kubernetesAdapter.buildOperatorTreeStructure(namespace);
        return ResponseEntity.ok(root);
    }

    @PostMapping("/operator/forward/{namespace}")
    @Operation(summary = "Forward operator tree structure to kirzop service",
              description = "Retrieves the operator-managed tree structure and forwards it to the kirzop service")
    public ResponseEntity<Void> forwardOperatorTreeStructure(
            @Parameter(description = "Kubernetes namespace to search for operator-managed tree pods")
            @PathVariable String namespace) {
        
        TreeNode root = kubernetesAdapter.buildOperatorTreeStructure(namespace);
        
        // Forward to kirzop service
        restTemplate.postForEntity("http://kirzop.default.svc.cluster.local/data", root, Void.class);
        
        return ResponseEntity.ok().build();
    }
} 
