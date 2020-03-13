FROM ethereum/client-go:stable

WORKDIR /gethnet
COPY datadir datadir/
COPY config.toml .

ENTRYPOINT [ "geth" ]
CMD [\
"--datadir", "/gethnet/datadir", \
"--mine", \
"--nodekey /gethnet/nodekey", \
"--ipcdisable", \
"--dev.period", "2", \
"--config", "/gethnet/config.toml" \
]
