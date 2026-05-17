# Build stage
FROM oven/bun:1-alpine AS builder

WORKDIR /app

# Copy package files
COPY package.json bun.lock ./

# Install dependencies
RUN bun install --frozen-lockfile

# Copy other frontend source files
COPY . .

# Build production bundle
RUN bun run build

# Production runtime stage
FROM oven/bun:1-alpine AS runner

WORKDIR /app

# Copy built server output
COPY --from=builder /app/.output ./.output

# Nitro server listens on port 3000 by default (can be overridden by PORT env)
EXPOSE 3000

ENV PORT=3000
ENV NODE_ENV=production

# Run Nuxt using Bun
CMD ["bun", ".output/server/index.mjs"]
