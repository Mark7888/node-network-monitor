import { create } from 'zustand';
import { Node, NodeDetails } from '../types/node.types';
import * as nodeService from '../services/nodeService';

interface NodesStore {
  nodes: Node[];
  selectedNode: NodeDetails | null;
  isLoading: boolean;
  error: string | null;
  
  fetchNodes: () => Promise<void>;
  fetchNodeDetails: (nodeId: string) => Promise<void>;
  clearError: () => void;
}

/**
 * Nodes store using Zustand
 */
export const useNodesStore = create<NodesStore>((set) => ({
  nodes: [],
  selectedNode: null,
  isLoading: false,
  error: null,

  fetchNodes: async () => {
    set({ isLoading: true, error: null });
    try {
      const data = await nodeService.getNodes();
      set({ nodes: data.nodes, isLoading: false });
    } catch (error: any) {
      set({
        error: error.response?.data?.error || 'Failed to fetch nodes',
        isLoading: false,
      });
    }
  },

  fetchNodeDetails: async (nodeId: string) => {
    set({ isLoading: true, error: null });
    try {
      const data = await nodeService.getNodeDetails(nodeId);
      set({ selectedNode: data, isLoading: false });
    } catch (error: any) {
      set({
        error: error.response?.data?.error || 'Failed to fetch node details',
        isLoading: false,
      });
    }
  },

  clearError: () => set({ error: null }),
}));
