import { create } from 'zustand';
import { Node, NodeDetails } from '../types/node.types';
import * as nodeService from '../services/nodeService';
import { showToast } from '@/shared/services/toastService';

interface NodesStore {
  nodes: Node[];
  selectedNode: NodeDetails | null;
  isLoading: boolean;
  error: string | null;
  
  fetchNodes: () => Promise<void>;
  fetchNodeDetails: (nodeId: string) => Promise<void>;
  archiveNode: (nodeId: string, archived: boolean) => Promise<void>;
  toggleFavorite: (nodeId: string) => Promise<void>;
  deleteNode: (nodeId: string) => Promise<void>;
  clearError: () => void;
}

/**
 * Nodes store using Zustand
 */
export const useNodesStore = create<NodesStore>((set, get) => ({
  nodes: [],
  selectedNode: null,
  isLoading: false,
  error: null,

  fetchNodes: async () => {
    set({ isLoading: true, error: null });
    try {
      const data = await nodeService.getNodes();
      set({ nodes: data.nodes || [], isLoading: false });
    } catch (error: unknown) {
      const err = error as { response?: { data?: { error?: string } } };
      set({
        error: err.response?.data?.error || 'Failed to fetch nodes',
        isLoading: false,
      });
    }
  },

  fetchNodeDetails: async (nodeId: string) => {
    set({ isLoading: true, error: null });
    try {
      const data = await nodeService.getNodeDetails(nodeId);
      set({ selectedNode: data, isLoading: false });
    } catch (error: unknown) {
      const err = error as { response?: { data?: { error?: string } } };
      set({
        error: err.response?.data?.error || 'Failed to fetch node details',
        isLoading: false,
      });
    }
  },

  archiveNode: async (nodeId: string, archived: boolean) => {
    try {
      await nodeService.archiveNode(nodeId, archived);
      
      // Update local state
      const nodes = get().nodes.map(node => 
        node.id === nodeId ? { ...node, archived } : node
      );
      set({ nodes });
      
      showToast.success(archived ? 'Node archived' : 'Node unarchived');
    } catch (error: unknown) {
      const err = error as { response?: { data?: { error?: string } } };
      const message = err.response?.data?.error || 'Failed to update node';
      showToast.error(message);
      throw error;
    }
  },

  toggleFavorite: async (nodeId: string) => {
    const node = get().nodes.find(n => n.id === nodeId);
    if (!node) return;

    const newFavorite = !node.favorite;
    
    try {
      await nodeService.setNodeFavorite(nodeId, newFavorite);
      
      // Update local state
      const nodes = get().nodes.map(n => 
        n.id === nodeId ? { ...n, favorite: newFavorite } : n
      );
      set({ nodes });
      
      showToast.success(newFavorite ? 'Added to favorites' : 'Removed from favorites');
    } catch (error: unknown) {
      const err = error as { response?: { data?: { error?: string } } };
      const message = err.response?.data?.error || 'Failed to update node';
      showToast.error(message);
      throw error;
    }
  },

  deleteNode: async (nodeId: string) => {
    try {
      await nodeService.deleteNode(nodeId);
      
      // Remove from local state
      const nodes = get().nodes.filter(node => node.id !== nodeId);
      set({ nodes });
      
      showToast.success('Node deleted successfully');
    } catch (error: unknown) {
      const err = error as { response?: { data?: { error?: string } } };
      const message = err.response?.data?.error || 'Failed to delete node';
      showToast.error(message);
      throw error;
    }
  },

  clearError: () => set({ error: null }),
}));
