FROM alpine:3.18.4

WORKDIR /app

COPY identity-server /app

# kubernetes will place additional config files in /app/config/ as volume
RUN ln -s /app/config/.env.local /app/.env.local

ENTRYPOINT ["/app/identity-server"]
