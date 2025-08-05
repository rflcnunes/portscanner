# Portscanner

Um scanner de portas TCP concorrente escrito em Go com interface de linha de comando moderna.

## Instalação

```bash
go install github.com/seu-usuario/portscanner@latest
```

## Funcionalidade do comando `scan`

O projeto define um executável chamado `portscanner` com um subcomando `scan`. A seguir, a descrição de cada parte para embasar a documentação:

### Estrutura geral

- **Root command (`portscanner`)**: ponto de entrada da CLI, exibe help geral.
- **Subcomando `scan`**: agrupa toda a lógica de escaneamento de portas.

### Flags do `scan`

| Flag      | Atalho | Tipo          | Padrão            | Descrição                                              |
| --------- | ------ | ------------- | ----------------- | ------------------------------------------------------ |
| --host    | -H     | string        | "scanme.nmap.org" | Host (IP ou domínio) onde executar o scan              |
| --start   | -s     | int           | 1                 | Porta inicial do intervalo a ser escaneado             |
| --end     | -e     | int           | 1024              | Porta final do intervalo a ser escaneado               |
| --timeout | -t     | time.Duration | 1s                | Tempo máximo de espera por conexão em cada tentativa   |
| --workers | -w     | int           | 100               | Número de goroutines concorrentes para dial simultâneo |

### Fluxo de execução (`runScan`)

1. **Leitura das flags**: Cobra faz o parse automático dos valores informados.
2. **Criação de canais**:
   - `ports chan int` bufferizado com capacidade igual a `workers`.
   - `results chan int` (sem buffer) para receber portas abertas.
3. **Inicialização dos workers**:
   - Cada worker é uma goroutine que lê de `ports`.
   - Para cada porta, executa `net.DialTimeout("tcp", host:porta, timeout)`.
   - Se conectar com sucesso, fecha a conexão e envia a porta para `results`.
4. **Enfileiramento de portas**:
   - Em outra goroutine, preenche `ports` de `start` até `end`, depois fecha o canal.
5. **Sincronização**:
   - Um `sync.WaitGroup` aguarda todos os workers terminarem antes de fechar `results`.
6. **Coleta e ordenação**:
   - Itera em `results`, acumula as portas em slice.
   - Ordena o slice com `sort.Ints`.
7. **Saída final**:
   - Exibe no console a lista de portas abertas e o tempo total de execução.

### Exemplos de uso

```sh
# Exibe ajuda
portscanner scan --help

# Exibe versão
portscanner --version

# Escaneia localhost de 1 a 1024 com 50 workers e timeout de 500ms
portscanner scan \
  --host localhost \
  --start 1 \
  --end 1024 \
  --workers 50 \
  --timeout 500ms
```

## Versionamento Automático

Este projeto utiliza GitHub Actions para versionamento automático:

- **Auto Version Bump**: Incrementa automaticamente a versão patch a cada push na branch `main`
- **Criação de Tags**: Cria tags Git automaticamente no formato `vX.Y.Z`
- **Releases**: Gera releases automaticamente no GitHub com notas de versão

### Como funciona

1. A cada push na branch `main`, o workflow extrai a versão atual do código
2. Incrementa o número do patch (ex: `1.0.0` → `1.0.1`)
3. Atualiza o código com a nova versão
4. Cria um commit e tag com a nova versão
5. Gera um release automaticamente no GitHub
