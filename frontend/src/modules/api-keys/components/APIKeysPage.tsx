import { useState, useEffect } from 'react';
import { useAPIKeys } from '../hooks/useAPIKeys';
import { CreateAPIKeyResponse } from '../types/apiKey.types';
import { formatTimestamp, formatRelativeTime } from '@/shared/utils/date';
import Card from '@/shared/components/ui/Card';
import Button from '@/shared/components/ui/Button';
import Badge from '@/shared/components/ui/Badge';
import Spinner from '@/shared/components/ui/Spinner';
import ErrorMessage from '@/shared/components/ui/ErrorMessage';
import EmptyState from '@/shared/components/ui/EmptyState';
import Modal from '@/shared/components/ui/Modal';
import Input from '@/shared/components/ui/Input';
import Select from '@/shared/components/ui/Select';
import { Plus, Key, Trash2, Copy, Check } from 'lucide-react';
import toast from 'react-hot-toast';

/**
 * API Keys management page component
 */
export default function APIKeysPage() {
  const { apiKeys, isLoading, error, fetchAPIKeys, createKey, toggleKey, deleteKey } = useAPIKeys();
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [isKeyDisplayModalOpen, setIsKeyDisplayModalOpen] = useState(false);
  const [newKeyName, setNewKeyName] = useState('');
  const [createdKey, setCreatedKey] = useState<CreateAPIKeyResponse | null>(null);
  const [copiedKey, setCopiedKey] = useState(false);
  const [isCreating, setIsCreating] = useState(false);
  const [filter, setFilter] = useState<'all' | 'enabled' | 'disabled'>('all');

  useEffect(() => {
    fetchAPIKeys();
  }, [fetchAPIKeys]);

  const handleCreateKey = async () => {
    if (!newKeyName.trim()) {
      toast.error('Please enter a name for the API key');
      return;
    }

    setIsCreating(true);
    const result = await createKey(newKeyName);
    setIsCreating(false);

    if (result) {
      setCreatedKey(result);
      setIsCreateModalOpen(false);
      setIsKeyDisplayModalOpen(true);
      setNewKeyName('');
    }
  };

  const handleCopyKey = async () => {
    if (createdKey?.key) {
      await navigator.clipboard.writeText(createdKey.key);
      setCopiedKey(true);
      toast.success('API key copied to clipboard');
      setTimeout(() => setCopiedKey(false), 2000);
    }
  };

  const handleDeleteKey = async (id: string) => {
    if (confirm('Are you sure you want to delete this API key? This action cannot be undone.')) {
      await deleteKey(id);
    }
  };

  // Filter API keys based on selected filter
  const filteredApiKeys = apiKeys.filter((apiKey) => {
    if (filter === 'enabled') return apiKey.enabled;
    if (filter === 'disabled') return !apiKey.enabled;
    return true; // 'all'
  });

  if (isLoading && apiKeys.length === 0) {
    return <Spinner message="Loading API keys..." />;
  }

  if (error) {
    return <ErrorMessage message={error} onRetry={fetchAPIKeys} />;
  }

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold">API Keys</h1>
          <p className="text-base-content/60 mt-1">
            Manage API keys for speedtest nodes
          </p>
        </div>
        <Button
          variant="primary"
          onClick={() => setIsCreateModalOpen(true)}
          className="gap-2"
        >
          <Plus size={18} />
          New Key
        </Button>
      </div>

      {/* Filter Dropdown */}
      {apiKeys.length > 0 && (
        <div className="flex justify-end">
          <div className="w-48">
            <Select
              value={filter}
              onChange={(e) => setFilter(e.target.value as 'all' | 'enabled' | 'disabled')}
              className="select-sm"
            >
              <option value="all">All Keys</option>
              <option value="enabled">Enabled Only</option>
              <option value="disabled">Disabled Only</option>
            </Select>
          </div>
        </div>
      )}

      {/* API Keys List */}
      {filteredApiKeys.length === 0 ? (
        <EmptyState
          title={apiKeys.length === 0 ? "No API keys found" : `No ${filter} API keys found`}
          message={
            apiKeys.length === 0
              ? "Create an API key to allow speedtest nodes to authenticate and submit measurements."
              : `There are no ${filter} API keys. Try changing the filter or ${filter === 'enabled' ? 'enable some keys' : 'disable some keys'}.`
          }
          icon={<Key size={64} />}
          action={
            apiKeys.length === 0 ? (
              <Button onClick={() => setIsCreateModalOpen(true)} variant="primary">
                Create First API Key
              </Button>
            ) : undefined
          }
        />
      ) : (
        <div className="space-y-4">
          {filteredApiKeys.map((apiKey) => (
            <Card key={apiKey.id} compact>
              <div className="flex justify-between items-start">
                <div className="flex-1">
                  <div className="flex items-center gap-3">
                    <h3 className="text-lg font-semibold">{apiKey.name}</h3>
                    <Badge variant={apiKey.enabled ? 'success' : 'ghost'}>
                      {apiKey.enabled ? 'Enabled' : 'Disabled'}
                    </Badge>
                  </div>
                  <p className="text-sm text-base-content/60 mt-1">
                    ID: {apiKey.id.substring(0, 16)}...
                  </p>
                  <div className="text-sm text-base-content/70 mt-2 space-y-1">
                    <p>Created: {formatTimestamp(apiKey.created_at)}</p>
                    <p>
                      Last used: {apiKey.last_used ? formatRelativeTime(apiKey.last_used) : 'Never'}
                    </p>
                  </div>
                </div>
                <div className="flex gap-2">
                  <Button
                    variant={apiKey.enabled ? 'ghost' : 'secondary'}
                    size="sm"
                    onClick={() => toggleKey(apiKey.id, !apiKey.enabled)}
                  >
                    {apiKey.enabled ? 'Disable' : 'Enable'}
                  </Button>
                  <Button
                    variant="error"
                    size="sm"
                    onClick={() => handleDeleteKey(apiKey.id)}
                    className="gap-1"
                  >
                    <Trash2 size={16} />
                    Delete
                  </Button>
                </div>
              </div>
            </Card>
          ))}
        </div>
      )}

      {/* Create API Key Modal */}
      <Modal
        isOpen={isCreateModalOpen}
        onClose={() => {
          setIsCreateModalOpen(false);
          setNewKeyName('');
        }}
        title="Create New API Key"
        actions={
          <>
            <Button variant="ghost" onClick={() => setIsCreateModalOpen(false)}>
              Cancel
            </Button>
            <Button
              variant="primary"
              onClick={handleCreateKey}
              loading={isCreating}
              disabled={isCreating}
            >
              Create
            </Button>
          </>
        }
      >
        <Input
          label="API Key Name"
          placeholder="e.g., Production Node 1"
          value={newKeyName}
          onChange={(e) => setNewKeyName(e.target.value)}
          helperText="Choose a descriptive name to identify this key"
          autoFocus
        />
      </Modal>

      {/* Display Created Key Modal */}
      <Modal
        isOpen={isKeyDisplayModalOpen}
        onClose={() => {
          setIsKeyDisplayModalOpen(false);
          setCreatedKey(null);
          setCopiedKey(false);
        }}
        title="API Key Created Successfully"
        size="md"
      >
        <div className="space-y-4">
          <div className="alert alert-warning">
            <div>
              <span className="font-semibold">⚠️ Important!</span>
              <p className="text-sm mt-1">
                Save this key securely. You won't be able to see it again.
              </p>
            </div>
          </div>

          <div>
            <label className="label">
              <span className="label-text font-semibold">Your API Key</span>
            </label>
            <div className="flex gap-2">
              <div className="flex-1 p-3 bg-base-200 rounded font-mono text-sm break-all">
                {createdKey?.key}
              </div>
              <Button
                variant={copiedKey ? 'secondary' : 'primary'}
                onClick={handleCopyKey}
                className="gap-2"
              >
                {copiedKey ? (
                  <>
                    <Check size={18} />
                    Copied
                  </>
                ) : (
                  <>
                    <Copy size={18} />
                    Copy
                  </>
                )}
              </Button>
            </div>
          </div>

          <div className="flex justify-end mt-6">
            <Button
              variant="primary"
              onClick={() => {
                setIsKeyDisplayModalOpen(false);
                setCreatedKey(null);
                setCopiedKey(false);
              }}
            >
              Close
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  );
}
