NAME = taskmaster

GO = go
GOFLAGS = -ldflags="-s -w"

all: $(NAME)

$(NAME):
	$(GO) build $(GOFLAGS) -o $(NAME) ./cmd/$(NAME)

clean:
	rm -f $(NAME)

fclean: clean
	$(GO) clean -cache

re: fclean all

run:
	go run ./cmd/taskmaster config.yml

test:
	$(GO) test ./...

lint:
	$(GO) vet ./...

.PHONY: all clean fclean re run test lint
