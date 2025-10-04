# Yarn Setup Instructions

This guide will help you install Yarn package manager for the DSM Parts Finder frontend.

## What is Yarn?

Yarn is a fast, reliable, and secure dependency management tool for JavaScript projects. It's an alternative to npm with improved performance and features.

## Installation Methods

### Method 1: Install via npm (Recommended)

If you already have Node.js and npm installed:

```bash
npm install -g yarn
```

### Method 2: Install via Corepack (Node.js 16.10+)

Corepack is included with Node.js 16.10+ and allows you to use Yarn without global installation:

```bash
corepack enable
corepack prepare yarn@stable --activate
```

### Method 3: Install via Package Manager

#### macOS (using Homebrew)
```bash
brew install yarn
```

#### Windows (using Chocolatey)
```bash
choco install yarn
```

#### Ubuntu/Debian
```bash
curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | sudo apt-key add -
echo "deb https://dl.yarnpkg.com/debian/ stable main" | sudo tee /etc/apt/sources.list.d/yarn.list
sudo apt update
sudo apt install yarn
```

#### CentOS/RHEL/Fedora
```bash
curl -sL https://dl.yarnpkg.com/rpm/yarn.repo -o /etc/yum.repos.d/yarn.repo
sudo yum install yarn
```

## Verify Installation

Check that Yarn is installed correctly:

```bash
yarn --version
```

You should see a version number like `1.22.19` or `3.6.4`.

## Quick Start for DSM Parts Finder

Once Yarn is installed, you can set up the frontend:

```bash
# Navigate to the frontend directory
cd frontend

# Install all dependencies
yarn install

# Start the development server
yarn dev
```

## Yarn vs npm Commands

Here's a quick reference for common commands:

| npm command | Yarn equivalent |
|-------------|----------------|
| `npm install` | `yarn` or `yarn install` |
| `npm install <package>` | `yarn add <package>` |
| `npm install --save-dev <package>` | `yarn add --dev <package>` |
| `npm uninstall <package>` | `yarn remove <package>` |
| `npm run <script>` | `yarn <script>` |
| `npm start` | `yarn start` |
| `npm test` | `yarn test` |
| `npm run build` | `yarn build` |

## Benefits of Using Yarn

- **Faster**: Parallel installation and caching
- **Reliable**: Uses lockfiles to ensure consistent installs
- **Secure**: Checksums to verify package integrity
- **Offline**: Can install packages from cache when offline
- **Workspaces**: Better support for monorepos

## Troubleshooting

### Common Issues

1. **Permission errors**: Use `sudo` on Unix systems or run as administrator on Windows
2. **Path issues**: Make sure Yarn is in your system PATH
3. **Version conflicts**: Use `yarn --version` to check your version

### Clear Cache

If you encounter issues, try clearing the Yarn cache:

```bash
yarn cache clean
```

### Reset Installation

To completely reset your dependencies:

```bash
rm -rf node_modules yarn.lock
yarn install
```

## Additional Resources

- [Official Yarn Documentation](https://yarnpkg.com/getting-started)
- [Yarn Migration Guide](https://yarnpkg.com/getting-started/migration)
- [Yarn CLI Commands](https://yarnpkg.com/cli)

## For DSM Parts Finder Development

After installing Yarn, you can use these project-specific commands:

```bash
# Install dependencies
yarn install

# Start development server (frontend will be at http://localhost:5173)
yarn dev

# Build for production
yarn build

# Preview production build
yarn preview

# Run linting
yarn lint
```

Make sure the Go backend is also running on port 8080 for full functionality.