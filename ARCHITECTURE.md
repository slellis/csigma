
```markdown
# ARCHITECTURE.md - Especificação Técnica do Compilador Sigma

Este documento registra as definições arquiteturais, escolhas de design e o protocolo de manutenção da linguagem **Sigma**.

## 1. Pipeline de Compilação
O Sigma opera em uma pipeline de cinco estágios distintos para garantir a separação de preocupações:

1.  **Source (.sig)**: Código fonte em texto puro.
2.  **Lexer**: Transformação de texto em tokens (usa mapa de keywords centralizado).
3.  **Parser**: Construção da Árvore de Sintaxe Abstrata (AST).
4.  **Semantic Analyzer**: Validação de regras de negócio e tipos (Fase Separada).
5.  **CodeGen**: Tradução da AST validada para **Assembly x86_64** (Linux).

---

## 2. Sistema de Tipos (Regra B: Tipagem Fixa)
O Sigma utiliza tipagem forte e estática. Uma vez que uma variável é declarada, seu tipo não pode ser alterado.

| Tipo Sigma | Tipo Go (Backend) | Descrição |
| :--- | :--- | :--- |
| SIGMA_INT | int64 | Inteiros de 64 bits. |
| SIGMA_FLT | float64 | Ponto flutuante de precisão dupla. |
| SIGMA_STR | string | Cadeias de caracteres. |
| SIGMA_BOOL | bool | Valores lógicos (true/false). |

---

## 3. Regras Semânticas e Operações

### 3.1 Aritmética e Divisão (Opção A)
* **Divisão Estrita**: O resultado da divisão segue o tipo dos operandos.
    * INT / INT = INT (Ex: 5 / 2 = 2)
    * FLT / FLT = FLT (Ex: 5.0 / 2.0 = 2.5)
* **Proibição de Coerção**: Operações entre tipos diferentes (ex: INT + FLT) resultam em erro semântico sem conversão explícita.

### 3.2 Conversão de Tipos (Casting)
A conversão deve ser sempre explícita, utilizando a sintaxe de função:
* float(expressão)
* int(expressão)
* str(expressão)

---

## 4. Gestão de Escopo e Símbolos
* **Escopo Global**: Para fins didáticos, todas as variáveis residem em um único escopo global.
* **Tabela de Símbolos**: Um mapa único armazena o par {Nome, Tipo}. Re-declarações do mesmo identificador no mesmo programa são proibidas.

```go
// Exemplo da estrutura interna (Em Go)
// Comentário didático: O tipo será fixado no momento da declaração (Regra B).
type Simbolo struct {
    Nome string
    Tipo string // SIGMA_INT, SIGMA_FLT, etc.
}

// TabelaSimbolos funciona como a "Fonte da Verdade" para o Analisador Semântico.
var TabelaSimbolos = make(map[string]Simbolo)

```

---

## 5. Tratamento de Erros

O Analisador Semântico é **resiliente**. Ele não interrompe a análise no primeiro erro encontrado.

* **Acumulador**: Todos os problemas são coletados em uma lista de strings (Slice).
* **Relatório**: O compilador exibe a lista completa de erros antes de abortar a fase de CodeGen.

---

## 6. Guia de Manutenção (Protocolo para Novos Comandos)

Para adicionar um novo comando (ex: if, while), siga esta ordem rigorosa:

1. **token/token.go**: Definir a constante do novo Token.
2. **lexer/lexer.go**: Adicionar a string no mapa keywords.
3. **parser/parser.go**: Implementar a lógica de construção do nó na AST.
4. **semantic/analyzer.go**: Adicionar a regra de validação e registro na Tabela de Símbolos.
5. **codegen/asm.go**: Definir a tradução para as instruções Assembly.

```

---
