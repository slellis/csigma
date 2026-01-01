Compilador CSigma (Versão Gold)

O CSigma é um compilador didático desenvolvido em Go que traduz código-fonte escrito na linguagem Sigma para Assembly x86_64, gerando executáveis binários reais para sistemas Linux.
🏛️ Filosofia do Projeto: Do Mainframe ao Registrador

Diferente de compiladores modernos que priorizam abstrações complexas, o CSigma foi concebido sob a ótica da disciplina de sistemas de grande porte (Mainframes).

Inspirado na organização rigorosa de sistemas clássicos, o CSigma separa claramente a intenção do programador:

    Data Division (Seção de Dados): Onde as variáveis são alocadas com precisão na memória.

    Procedure Division (Seção de Código): Onde a lógica flui de forma linear, gerando um Assembly limpo, alinhado e 100% comentado.

🛠️ O Coração do Compilador (Explicando o Go)

Para garantir a transparência do processo, o CSigma utiliza recursos estratégicos da linguagem Go. Abaixo, detalhamos algumas escolhas técnicas cruciais:

    Manipulação de Arquivos e Sufixos: No arquivo main.go, utilizamos a lógica strings.TrimSuffix(inputPath, ".sig") + ".log".

        strings.TrimSuffix: Esta função identifica o nome do arquivo fonte e remove a extensão original .sig.

        + ".log": Acrescentamos o novo sufixo para garantir que cada compilação gere um rastro técnico (log) único com o mesmo nome do programa.

    A Estratégia io.MultiWriter: Implementamos o MultiWriter para o modo Verbose. Isso permite que o compilador envie dados simultaneamente para o terminal (os.Stdout) e para o arquivo de log, garantindo que o rastro da compilação seja registrado permanentemente.

    Diferenciação de Operandos no Codegen: O gerador de código identifica se um valor é um Literal (número puro) ou um Identificador (variável). Isso decide se o Assembly gerado será um mov rax, 100 (valor imediato) ou mov rax, [A] (busca em memória), garantindo a integridade da execução e evitando falhas de proteção de memória.

🚀 O Pipeline de Compilação

O CSigma percorre quatro fases distintas até entregar o binário final:

    Análise Léxica (Lexer): Escaneia o texto fonte e gera Tokens (unidades básicas).

    Análise Sintática (Parser): Constrói a AST (Abstract Syntax Tree), que é o mapa lógico e hierárquico das instruções.

    Geração de Código (Codegen): Traduz a AST para instruções Assembly x86_64 devidamente comentadas.

    Montagem e Linkagem: Utiliza o NASM (Assembler) e o GCC (Linker) para criar o executável final.

📝 Exemplo de Código Sigma

Abaixo, um exemplo de uma calculadora interativa que demonstra a capacidade atual da linguagem:
Snippet de código

// Declaração de Variáveis
VAR A 0
VAR B 0
VAR C 0
VAR RESULTADO 0

// Entrada de Dados
PRINT "VALOR A:"
INPUT A
PRINT "VALOR B:"
INPUT B
PRINT "VALOR C:"
INPUT C

// Processamento Aritmético (Expressão Complexa)
RESULTADO = A + B * 2 / C

// Saída dos Resultados
PRINT "RESULTADO FINAL:"
PRINT RESULTADO

📊 Relatório de LOG (Listing de Compilação)

Ao compilar, o CSigma gera um arquivo .log detalhado que funciona como um "Listing" de Mainframe, contendo:

    Trace de Tokens: Cada unidade identificada pelo Lexer com seu tipo e conteúdo.

    Dump da AST: A representação estrutural da árvore sintática para conferência lógica.

    Status de Build: O passo a passo das chamadas externas ao NASM e GCC.

⚙️ Pré-requisitos e Execução

Para rodar este compilador, você precisará de:

    Go (v1.18 ou superior)

    NASM (Netwide Assembler)

    GCC (GNU Compiler Collection)

Como Compilar e Rodar:
Bash

# Executa o compilador passando o arquivo Sigma
go run main.go exemplos/calculadora.sig

# Executa o binário gerado
./calculadora

Desenvolvido por Sidney Unindo a experiência dos sistemas de grande porte com a agilidade do desenvolvimento moderno.
