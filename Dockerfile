# LICENSE UPL 1.0
#
# Copyright (c) 2014, 2019, Oracle and/or its affiliates. All rights reserved.
#
# ORACLE DOCKERFILES PROJECT
# --------------------------
#
# Dockerfile template for Oracle Instant Client
#
# HOW TO BUILD THIS IMAGE AND RUN A CONTAINER
# -------------------------------------------
#
# Run:
#      $ docker build --pull -t oracle/instantclient:19 .
#      $ docker run -ti --rm oracle/instantclient:19 sqlplus hr@example.com/orclpdb1
#
# NOTES
# -----
#
# Applications using Oracle Call Interface (OCI) 19 can connect to
# Oracle Database 11.2 or later.  Some tools may have other
# restrictions.
#
# Note Instant Client 19 automatically configures the global library
# search path to include Instant Client libraries.
#
# OPTIONAL ORACLE CONFIGURATION FILES
# -----------------------------------
#
# Optional Oracle Network and Oracle client configuration files can be put in the
# default configuration file directory /usr/lib/oracle/<version>/client64/lib/network/admin.
# Configuration files include tnsnames.ora, sqlnet.ora, oraaccess.xml and
# cwallet.sso.  You can use a Docker volume to mount the directory containing
# the files at runtime, for example:
#
#   docker run -v /my/host/wallet_dir:/usr/lib/oracle/19.3/client64/lib/network/admin:Z,ro . . .
#
# This avoids embedding private information such as wallets in images.  If you
# do choose to include network configuration files in images, you can use a
# Dockerfile COPY, for example:
#
#   COPY tnsnames.ora sqlnet.ora /usr/lib/oracle/${release}.${update}/client64/lib/network/admin/
#
# There is no need to set the TNS_ADMIN environment variable when files are in
# the container's default configuration file directory, as shown.
#
# ORACLE INSTANT CLIENT PACKAGES
# ------------------------------
#
# Instant Client Packages are available from https://yum.oracle.com/repo/OracleLinux/OL7/oracle/instantclient/x86_64/index.html
#
# Base - one of these packages is required to run applications and tools
#   oracle-instantclientXX.Y-basic      : Basic Package - All files required to run OCI, OCCI, and JDBC-OCI applications
#   oracle-instantclientXX.Y-basiclite  : Basic Light Package - Smaller version of the Basic package, with only English error messages and Unicode, ASCII, and Western European character set support
#
# Tools - optional packages (requires the 'basic' package)
#   oracle-instantclientXX.Y-sqlplus    : SQL*Plus Package - The SQL*Plus command line tool for SQL and PL/SQL queries
#   oracle-instantclientXX.Y-tools      : Tools Package - Includes Data Pump, SQL*Loader and Workload Replay Client
#
# Development and Runtime - optional packages (requires the 'basic' package)
#   oracle-instantclientXX.Y-devel      : SDK Package - Additional header files and an example makefile for developing Oracle applications with Instant Client
#   oracle-instantclientXX.Y-jdbc       : JDBC Supplement Package - Additional support for Internationalization under JDBC
#   oracle-instantclientXX.Y-odbc       : ODBC Package - Additional libraries for enabling ODBC applications
#

FROM oraclelinux:8-slim AS instantclient

RUN  microdnf -y update && microdnf install -y oracle-instantclient-release-el8 && \
     microdnf install oracle-instantclient-basic && \
     rm -rf /var/cache/dnf

# Uncomment if the tools package is added
# ENV PATH=$PATH:/usr/lib/oracle/${release}.${update}/client64/bin

# Start from a Instant Client Image
FROM instantclient AS builder

# Add Maintainer Info
LABEL maintainer="Rui Miranda <a75488@alunos.uminho.pt>"

# Phase 1: Go microservice compilation

## Install Golang + git
RUN microdnf -y module enable go-toolset:ol8 && \
    microdnf -y install go-toolset && \
    rm -rf /var/cache/dnf

# Set the Current Working Directory inside the container
WORKDIR /mcdt

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# Phase 2: Just Instant Client + compiled application

FROM instantclient

# Copy compiled file to final container
COPY --from=builder /mcdt/main .

COPY --from=builder /mcdt/swagger ./swagger

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
