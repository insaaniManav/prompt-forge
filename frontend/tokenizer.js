// Token Counter using js-tiktoken (local package)
class TokenCounter {
    constructor() {
        this.encoding = null;
        this.callbacks = [];
        this.currentStats = { tokens: 0 };
        this.initializeTiktoken();
    }

    async initializeTiktoken() {
        try {
            // Try to load tiktoken if available (browser environment doesn't support require)
            if (typeof window !== 'undefined' && window.tiktoken) {
                this.encoding = window.tiktoken.encoding_for_model('gpt-4');
                console.log('✅ Tiktoken initialized from window object');
            } else {
                console.log('ℹ️ Tiktoken not available in browser, using estimation');
                this.encoding = null;
            }
            
            // Update count if there's already text
            const promptTextarea = document.getElementById('main-prompt');
            if (promptTextarea && promptTextarea.value) {
                this.updateCount(promptTextarea.value);
            }
        } catch (error) {
            console.warn('⚠️ Tiktoken failed to load, falling back to simple counting:', error);
            this.encoding = null;
        }
    }

    // Register callback for token count updates
    onUpdate(callback) {
        this.callbacks.push(callback);
    }

    // Update token count for given text
    updateCount(text) {
        let tokenCount;
        
        if (this.encoding) {
            try {
                const tokens = this.encoding.encode(text);
                tokenCount = tokens.length;
            } catch (error) {
                console.warn('Token encoding error:', error);
                tokenCount = this.estimateTokens(text);
            }
        } else {
            tokenCount = this.estimateTokens(text);
        }

        this.currentStats = { tokens: tokenCount };
        
        // Notify all callbacks
        this.callbacks.forEach(callback => {
            try {
                callback(this.currentStats);
            } catch (error) {
                console.error('Token counter callback error:', error);
            }
        });
    }

    // Fallback token estimation (roughly 4 characters per token)
    estimateTokens(text) {
        if (!text) return 0;
        
        // More sophisticated estimation
        const words = text.trim().split(/\s+/).length;
        const characters = text.length;
        
        // Average estimation: ~0.75 tokens per word, but at least chars/4
        return Math.max(Math.ceil(words * 0.75), Math.ceil(characters / 4));
    }

    // Get current token count
    getCurrentCount() {
        return this.currentStats.tokens;
    }
}

// Model context information
function getModelContextInfo() {
    // Get all models from provider configurations
    const allModels = {
        // OpenAI models
        'gpt-4': { name: 'GPT-4', limit: 8192, formattedLimit: '8K' },
        'gpt-4-turbo': { name: 'GPT-4 Turbo', limit: 128000, formattedLimit: '128K' },
        'gpt-3.5-turbo': { name: 'GPT-3.5 Turbo', limit: 16384, formattedLimit: '16K' },
        'o4-mini': { name: 'O4 Mini', limit: 128000, formattedLimit: '128K' },
        
        // Azure OpenAI models
        'gpt-4.1': { name: 'GPT-4.1', limit: 200000, formattedLimit: '200K' },
        'o3': { name: 'O3', limit: 1000000, formattedLimit: '1M' },
        
        // Anthropic Claude models
        'claude-3-5-sonnet-20241022': { name: 'Claude 3.5 Sonnet', limit: 200000, formattedLimit: '200K' },
        'claude-3-haiku-20240307': { name: 'Claude 3 Haiku', limit: 200000, formattedLimit: '200K' },
        'claude-3-opus-20240229': { name: 'Claude 3 Opus', limit: 200000, formattedLimit: '200K' },
        
        // Ollama models
        'gemma3:12b': { name: 'Gemma 3 12B', limit: 128000, formattedLimit: '128K' },
        'llama3.2:3b': { name: 'Llama 3.2 3B', limit: 128000, formattedLimit: '128K' },
        'llama3.1:8b': { name: 'Llama 3.1 8B', limit: 128000, formattedLimit: '128K' },
        'mistral:7b': { name: 'Mistral 7B', limit: 32000, formattedLimit: '32K' },
        'codellama:7b': { name: 'Code Llama 7B', limit: 16000, formattedLimit: '16K' },
        'phi3:3.8b': { name: 'Phi-3 3.8B', limit: 128000, formattedLimit: '128K' },
        'gemma2:2b': { name: 'Gemma 2 2B', limit: 8000, formattedLimit: '8K' },
        'gemma2:9b': { name: 'Gemma 2 9B', limit: 8000, formattedLimit: '8K' }
    };
    
    // Get current model from selection or default to Gemma 3 12B (Ollama default)
    const currentModel = (typeof getCurrentModel === 'function') ? getCurrentModel() : 'gemma3:12b';
    return allModels[currentModel] || allModels['gemma3:12b'];
}

// Token warning system
function getTokenWarning(tokenCount) {
    const modelInfo = getModelContextInfo();
    const warnings = [];
    
    const warningThreshold = modelInfo.limit * 0.8; // 80% of limit
    const dangerThreshold = modelInfo.limit * 0.95; // 95% of limit
    
    if (tokenCount > dangerThreshold) {
        warnings.push({
            status: 'danger',
            message: `⚠️ Very close to ${modelInfo.formattedLimit} limit`,
            tokens: `${tokenCount.toLocaleString()}/${modelInfo.limit.toLocaleString()}`
        });
    } else if (tokenCount > warningThreshold) {
        warnings.push({
            status: 'warning',
            message: `⚡ Approaching ${modelInfo.formattedLimit} limit`,
            tokens: `${tokenCount.toLocaleString()}/${modelInfo.limit.toLocaleString()}`
        });
    }
    
    return warnings;
}

// Initialize token counter function
function initializeTokenCounter() {
    window.tokenCounter = new TokenCounter();
    console.log('🔢 Token counter initialized');
}

// Export for use in other files
if (typeof module !== 'undefined' && module.exports) {
    module.exports = { TokenCounter, getModelContextInfo, getTokenWarning, initializeTokenCounter };
} 