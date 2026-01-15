# CSigma Compiler - Platinum Edition 🚀

O **CSigma** é um compilador de 64 bits desenvolvido em Go, projetado para traduzir a linguagem Sigma diretamente para **Assembly x86_64 (NASM)**, com posterior linkagem via **GCC**. 

O projeto demonstra as etapas fundamentais da construção de um compilador: análise léxica, sintática, geração de código e integração com bibliotecas de baixo nível (LibC).



## 🛠️ Status da Versão: Platinum
Atualmente, o compilador é capaz de processar aritmética linear, realizar entrada e saída de dados via terminal e gerar binários executáveis reais.

### Funcionalidades Atuais:
* **Aritmética Linear:** Suporte para as quatro operações básicas (`+`, `-`, `*`, `/`) em expressões encadeadas.
* **Interatividade (I/O):** Implementação dos comandos `print` (para strings e variáveis) e `input` (para captura de dados via teclado).
* **Integração com LibC:** O código gerado utiliza as funções `printf` e `scanf` da biblioteca padrão do C.
* **Relatório Técnico (Verbose Mode):** Geração automática de Logs detalhados com Dump da **AST (Abstract Syntax Tree)**, listagem de Tokens e o código Assembly final.
* **Target x86_64:** Geração de código Assembly NASM puro para Linux 64 bits.

---

## 🏗️ Arquitetura do Sistema

1.  **Lexer (Scanner):** Converte o código fonte em tokens lógicos. Suporta comentários de linha (`//`), strings e números decimais.
2.  **Parser (Analista Sintático):** Reconhece a gramática e constrói a **AST** via *Recursive Descent*.
3.  **CodeGen (Gerador de Código):** Traduz a AST para x86_64, gerenciando registradores (`RAX`, `RBX`, `RDI`, `RSI`) e alinhamento de pilha.
4.  **Linker (GCC):** Realiza a montagem e linkagem final com a LibC.



---

## 🚀 Como Executar

### Pré-requisitos:
* **Go** (1.18+)
* **NASM** (Assembler)
* **GCC** (Linker)

### Compilando um código Sigma:
```bash
# Execute o compilador passando seu código fonte
go run main.go exemplos/calculadora.sig

# O compilador gerará o executável com o nome do arquivo fonte:
./calculadora

📊 Exemplo de Código Sigma
Snippet de código

// TESTE DAS OPERACOES NO CSIGMA
print "Calculadora Platinum"

var a = 0
var b = 0
var res = 0

print "Digite o valor de a:"
input a
print "Digite o valor de b:"
input b

res = a + b * 2
print "Resultado final:"
print res

🗺️ Roadmap: Rumo à Versão Diamond

    [ ] Reativação do Semantic Analyzer: Validação de tipos e escopo.

    [ ] Estruturas de Controle: Implementação de IF e FOR.

    [ ] Precedência Matemática: Suporte a parênteses () e ordem de operações.

Desenvolvido por: Sidney (2026)


---

### 2. Arquivo `.gitignore` (Obrigatório para um bom repositório)
Crie um arquivo chamado `.gitignore` na raiz do projeto e coloque isso dentro. Isso impedirá que arquivos temporários de compilação sejam enviados para o seu GitHub.

```text
# Binários e Objetos
*.o
*.out
output.asm

# Executáveis gerados (nomes comuns)
calculadora
programa
teste

# Logs de compilação
*.log

# Binários do Go
csigma