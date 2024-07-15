FROM node:20

# Create and change to the app directory.
WORKDIR /usr/src/app

# Copy package.json, package-lock.json, and hardhat config files
# Assuming hardhat config files are in the root of the hardhat directory
COPY ./hardhat/package*.json ./
COPY ./hardhat/hardhat.config.ts ./

# Install dependencies
RUN npm install

# Copy the rest of your hardhat application
COPY ./hardhat/* ./

# Expose the port the app runs on
EXPOSE 8545

# Run the chain command as defined in package.json
CMD ["npm", "run", "chain"]
