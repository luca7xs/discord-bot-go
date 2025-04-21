# Discord Bot em Go

Um bot para Discord desenvolvido em Go, com suporte a comandos, criação e fechamento de tickets, verificação de administradores, sistema de logs e mais!

---

## 📁 Estrutura do Projeto

```
cmd/
└── main.go               # Ponto de entrada da aplicação

internal/
├── bot/
│   └── bot.go            # Inicialização do bot
│
├── commands/             # Comandos de interação do bot
│   ├── avatar.go
│   ├── hello.go
│   ├── info.go
│   ├── ping.go
│   ├── registry.go
│   └── ticket.go
│
├── components/           # Componentes interativos (buttons, modals)
│   ├── close_ticket.go
│   ├── confirm_close.go  # Confirmação de fechamento com ou sem logs
│   ├── create_ticket.go
│   ├── registry.go
│   ├── ticket_reason.go
│   └── ticket_type.go
│
├── config/
│   └── config.go         # Configurações do bot (tokens, etc)

├── db/                   # Conexão com banco de dados e models
│   ├── connection.go
│   ├── models.go
│   └── tickets.go
│
└── utils/                # Funções utilitárias
    ├── admin_check.go
    ├── fetch_messages.go
    └── response.go
```

---

## ⚙️ Como rodar

### 1. Clone o repositório

```bash
git clone https://github.com/luca7xs/discord-bot-go.git
cd discord-bot-go
```

### 2. Instale as dependências

```bash
go mod tidy
```

### 3. Configure o `config.go` em `internal/config/config.go`

Preencha os dados do seu bot, como token e outras informações necessárias.

### 4. Execute o bot

```bash
go run cmd/main.go
```

---

## 📌 Funcionalidades

- ✅ Comando `/ping`, `/info`, etc.
- 🎟️ Sistema de tickets com criação, motivo e encerramento
- 🗂️ Salvamento de logs dos tickets (mensagens)
- ❌ Possibilidade de cancelar o fechamento do ticket
- 🔒 Verificação se o usuário é admin
- 🧠 Registro automático dos comandos no Discord
- 💾 Integração com banco de dados (MySQL)

---

## 🧾 Sistema de Logs de Tickets

Durante o processo de fechamento de um ticket, o bot realiza:

1. **Pergunta se deseja mesmo fechar o ticket**  
2. **Pergunta se deseja salvar os logs**

Se confirmado, ele:

✅ Busca todas as mensagens do canal  
✅ Salva essas mensagens no banco de dados via `db.CloseTicket`  
✅ Exclui o canal após 20 segundos (com opção de cancelar via botão)

**Local do código:**  
`internal/components/confirm_close.go`

**Exemplo de mensagem no canal:**

> O ticket de **[tipo]** será fechado com logs salvos!  
> Esse canal será excluído em 20 segundos.  
> *(Botão de Cancelar disponível)*

---

## 🛠️ Tecnologias

- [Go](https://golang.org)
- [DiscordGo](https://github.com/bwmarrin/discordgo)
- Banco de dados (configurado em `db/connection.go`)

---

## 🧑‍💻 Contribuindo

Contribuições são bem-vindas! Sinta-se livre para abrir issues e PRs.

---

Feito com 💙 por [luca7xs](https://github.com/luca7xs)
