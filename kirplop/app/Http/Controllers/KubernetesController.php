<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Http;

class KubernetesController extends Controller
{
    public function show()
    {
        return view('kubernetes');
    }

    public function submit(Request $request)
    {
        $request->validate([
            'number' => 'required|numeric|min:1',
            'name' => 'required|regex:/^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/'
        ]);

        try {
            // Get the values from the request
            $number = $request->input('number');
            $name = $request->input('name');

            $apiUrl = env('KUBERNETES_API_URL') . '/apis/tree.nullcloud.io/v1alpha1/namespaces/default/trees';
            $headers = [
                'Authorization' => 'Bearer ' . env('KUBERNETES_API_TOKEN'),
            ];

            // First, check if the resource exists
            $getResponse = Http::withHeaders($headers)
                ->withOptions(['verify' => false])
                ->get($apiUrl . '/' . $name);

            if ($getResponse->successful()) {
                // Get the existing resource data
                $existingResource = $getResponse->json();
                
                // Prepare the payload for update, including the resourceVersion
                $payload = [
                    'apiVersion' => 'tree.nullcloud.io/v1alpha1',
                    'kind' => 'Tree',
                    'metadata' => [
                        'name' => $name,
                        'resourceVersion' => $existingResource['metadata']['resourceVersion'],
                        'labels' => [
                            'app.kubernetes.io/name' => 'tree-operator',
                            'app.kubernetes.io/managed-by' => 'kustomize'
                        ]
                    ],
                    'spec' => [
                        'count' => (int)$number
                    ]
                ];

                // Resource exists, update it
                $response = Http::withHeaders($headers)
                    ->withOptions(['verify' => false])
                    ->put($apiUrl . '/' . $name, $payload);
                
                $message = 'Resource ' . $name . ' was successfully updated with ' . $number . ' vertices!';
            } else {
                // Prepare the payload for creation
                $payload = [
                    'apiVersion' => 'tree.nullcloud.io/v1alpha1',
                    'kind' => 'Tree',
                    'metadata' => [
                        'name' => $name,
                        'labels' => [
                            'app.kubernetes.io/name' => 'tree-operator',
                            'app.kubernetes.io/managed-by' => 'kustomize'
                        ]
                    ],
                    'spec' => [
                        'count' => (int)$number
                    ]
                ];

                // Resource doesn't exist, create it
                $response = Http::withHeaders($headers)
                    ->withOptions(['verify' => false])
                    ->post($apiUrl, $payload);
                
                $message = 'Resource ' . $name . ' was successfully created with ' . $number . ' vertices!';
            }

            if ($response->successful()) {
                return redirect()->back()->with('success', $message);
            } else {
                return redirect()->back()->with('error', 'Failed to communicate with Kubernetes: ' . $response->body());
            }
        } catch (\Exception $e) {
            return redirect()->back()->with('error', 'An error occurred: ' . $e->getMessage());
        }
    }
} 