FROM node:20-alpine

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies for building
RUN npm ci

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Install serve globally for production static file serving
RUN npm install -g serve

# Copy entrypoint script
COPY docker-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:3000/ || exit 1

# Use entrypoint script to generate runtime config
ENTRYPOINT ["docker-entrypoint.sh"]
