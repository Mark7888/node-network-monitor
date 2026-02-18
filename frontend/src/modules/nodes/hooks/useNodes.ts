import { useNodesStore } from '../store/nodesStore';

/**
 * Custom hook for nodes
 */
export function useNodes() {
  const {
    nodes,
    selectedNode,
    isLoading,
    error,
    fetchNodes,
    fetchNodeDetails,
    archiveNode,
    toggleFavorite,
    deleteNode,
    clearError,
  } = useNodesStore();

  return {
    nodes,
    selectedNode,
    isLoading,
    error,
    fetchNodes,
    fetchNodeDetails,
    archiveNode,
    toggleFavorite,
    deleteNode,
    clearError,
  };
}
