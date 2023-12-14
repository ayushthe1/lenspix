FROM node:latest

RUN mkdir /tailwind

# Set /tailwind as the workdir.
# A workdir is required for npm to work correctly.
WORKDIR /tailwind

# Install tailwindcss and initialize
RUN npm init -y && \
    npm install tailwindcss && \
    npx tailwindcss init

# Run tailwindcss. This will watch for changes in /src/styles.css and output to /dst/styles.css
CMD npx tailwindcss -c /src/tailwind.config.js -i /src/styles.css -o /dst/styles.css --watch --poll