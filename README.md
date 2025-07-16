# PromptForge 🔨

**The AI prompt engineering workbench that generates, analyzes, and systematically tests prompts**

Stop writing prompts from scratch. Let AI help you craft better prompts, then systematically validate them. Built with Go for speed and reliability.

**[🎮 Try Interactive Demo](https://demo.arcade.software/bFGTYb7AuRV33Kei7ZFQ)**

![PromptForge Screenshot](screenshot.png)

## Why PromptForge?

Most prompt tools are glorified text editors. PromptForge brings **engineering discipline** to prompt development:

- **✨ Smart Prompt Generation** - AI-powered prompt crafting with intelligent suggestions
- **🔍 Advanced Prompt Analysis** - Get optimization feedback before you even test
- **🧪 Systematic Evaluation** - Generate comprehensive test suites automatically
- **📊 Performance Analytics** - Track what works across different scenarios
- **📚 Version Control** - Never lose a working prompt again

## Quick Start

**One-line Docker setup:**
```bash
docker run -d -p 8080:8080 -e ANTHROPIC_API_KEY="your-key" ghcr.io/insaanimanav/prompt-forge:main
```

**Or clone and run:**
```bash
git clone https://github.com/insaanimanav/promptforge.git
cd promptforge && ./start.sh
```

Open `http://localhost:8080` and start crafting better prompts.

## Core Features

### ✨ AI-Powered Prompt Generation
- **Smart Crafting**: AI assists in building effective prompts from scratch
- **Intelligent Suggestions**: Context-aware recommendations for prompt improvement

### 🔍 Advanced Prompt Analysis
- **Quick Review**: Instant optimization suggestions
- **Deep Analysis**: Comprehensive prompt evaluation and scoring
- **Best Practice Validation**: Checks against proven prompt engineering principles

### 🧪 Systematic Evaluation Engine
Creates comprehensive test suites including:
- **Robustness**: Edge cases, typos, input variations
- **Safety**: Bias detection, harmful content resistance
- **Accuracy**: Factual correctness verification
- **Creativity**: Novel thinking scenarios

### ⚡ Professional Testing
- Execute with full parameter control
- Dynamic variable detection and substitution
- Multi-model comparison side by side (Claude, GPT-4, Azure OpenAI, Ollama)

### 📚 Prompt Management
- Organized library with search and tags
- Complete execution history with filtering
- Export/import capabilities

## Supported Models

- **Claude 3.5 Sonnet** (200K context) - Excellent reasoning
- **GPT-4.1** (200K context) - Detailed analysis  
- **O3** (1M context) - Fast execution
- **Azure OpenAI** - Enterprise-ready
- **Ollama Models** - Local AI with Llama, Mistral, Code Llama, and more

## Configuration

### Docker (Recommended)
```bash
# Anthropic (Default)
-e ANTHROPIC_API_KEY="sk-ant-api03-..."

# OpenAI  
-e OPENAI_API_KEY="sk-..."

# Azure OpenAI
-e AZURE_OPENAI_API_KEY="your-key"
-e AZURE_OPENAI_BASE_URL="https://your-resource.openai.azure.com"

# Ollama (Local AI)
-e OLLAMA_BASE_URL="http://localhost:11434"
-e DEFAULT_AI_PROVIDER="ollama"
```

### Local Development
```bash
# Clone and setup
git clone https://github.com/insaanimanav/promptforge.git
cd promptforge

# Configure environment
cp .env.example .env
# Add your API keys to .env

# Run
cd api && go run main.go
```

## API Reference

Core endpoints for integration:
- `POST /api/critique` - Analyze prompts
- `POST /api/execute` - Test with parameters
- `POST /api/generate-eval` - Create test suites
- `GET /api/prompts` - Manage prompt library

## Contributing

Built for the prompt engineering community. Contributions welcome:

- 🌐 Additional AI model integrations
- 📊 Enhanced evaluation metrics
- 🎨 UI/UX improvements
- 🔧 Performance optimizations

1. Fork the repository
2. Create feature branch: `git checkout -b feature/name`
3. Submit pull request

## License

GPLv3 - Built with ❤️ for better AI interactions

---

**Transform prompt engineering from art to science**