# Scanscope

Scanscope é uma ferramenta de varredura de código-fonte configurável em Go. As regras são definidas em `config/rules.json` e podem conter etapas encadeadas e validação opcional com a OpenAI.

## Recursos

- Regras hierárquicas em JSON
- Cache baseado em hash SHA-256 do arquivo e categoria
- Suporte a múltiplas extensões de arquivos
- Validação semântica opcional via IA
- Relatórios em `report.json` e `report.md`

## Execução

```bash
# variáveis como OPENAI_API_KEY devem estar definidas para uso da IA
 go run ./cmd/scanner/main.go
```
