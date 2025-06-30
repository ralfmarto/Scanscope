# Scanscope

**Scanscope** é uma ferramenta de varredura de código-fonte altamente configurável, que combina expressões regulares, validações encadeadas e inteligência artificial para identificar riscos reais com mais precisão e menos falsos positivos.

## ✨ Principais recursos

- ✅ Análise estática por regras definidas em JSON
- 🔁 Encadeamento de regex e validações por contexto
- 🤖 Integração com IA (OpenAI) para validação semântica
- 📁 Suporte a múltiplas extensões de arquivos
- 🧠 Cache inteligente por hash + categoria
- 📊 Geração de relatórios em JSON e Markdown
- ⚙️ Integração com GitHub Actions para uso em CI/CD

## 💼 Casos de uso

- Detecção de dados sensíveis (CPF, CNPJ, senhas)
- Verificação de logs mal utilizados
- Análise personalizada de regras internas de SAST
- Criação de políticas de segurança como código

## 🚀 Como executar

```bash
go run ./cmd/scanner/main.go
