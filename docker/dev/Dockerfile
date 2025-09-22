# Development container for Accumulate DID project
# Provides Go 1.23, Node.js LTS, and all necessary build tools

FROM golang:1.23-bookworm

# Install system dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    git \
    curl \
    make \
    jq \
    unzip \
    ca-certificates \
    gnupg \
    && rm -rf /var/lib/apt/lists/*

# Install Node.js LTS via NodeSource
RUN curl -fsSL https://deb.nodesource.com/setup_lts.x | bash - \
    && apt-get install -y nodejs \
    && rm -rf /var/lib/apt/lists/* \
    && node -v && npm -v

# Create non-root user 'dev' with uid 1000
RUN groupadd -g 1000 dev \
    && useradd -m -u 1000 -g dev -s /bin/bash dev

# Set working directory
WORKDIR /workspace

# Set environment variables for Go
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOPROXY=https://proxy.golang.org,direct \
    GOSUMDB=sum.golang.org

# Verify installations
RUN go version \
    && node --version \
    && npm --version \
    && git --version

# Switch to non-root user
USER dev

# Set up bash prompt with colors and info - need to set PATH for Go
RUN echo 'export PATH=/usr/local/go/bin:$PATH' >> ~/.bashrc \
    && echo 'export PS1="\[\033[36m\]\u@dev-container\[\033[m\]:\[\033[33;1m\]\w\[\033[m\]\$ "' >> ~/.bashrc \
    && echo 'alias ll="ls -la"' >> ~/.bashrc \
    && echo 'alias gs="git status"' >> ~/.bashrc \
    && echo 'echo "🚀 Accumulate DID Development Container"' >> ~/.bashrc \
    && echo 'echo "   Go: $(go version | cut -d" " -f3)"' >> ~/.bashrc \
    && echo 'echo "   Node: $(node -v)"' >> ~/.bashrc \
    && echo 'echo "   npm: $(npm -v)"' >> ~/.bashrc

# Default command
CMD ["/bin/bash"]