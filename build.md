docker build -f Dockerfile.alpine .  -t npolizotis/sftpgo:2 --build-arg FEATURES=nos3,nogcs --build-arg INSTALL_OPTIONAL_PACKAGES=true

docker run --name sftp2 -p 10022:2022 -p 8080:8080 -p 10023:2023 -p 10024:2024 npolizotis/sftpgo:2

# test