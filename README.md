# kubecompare

A CLI tool to compare different revisions of Kubernetes resources.

## Description

`kubecompare` is a command-line utility written in Go that helps users compare different revisions of Kubernetes deployments, daemonsets, or statefulsets to understand what changed between versions.

## Installation

### From Source

Clone the repository and build the application using Go:

```bash
git clone https://github.com/your-username/kubecompare.git
cd kubecompare
go build -o kubecompare
```

## Homebrew

If you have Homebrew installed, you can install kubecompare from our custom tap:

```
brew tap your-username/my-tap
brew install kubecompare
```

## Usage

```
kubecompare [<resource-type> <resource-name> | <resource-type>/<resource-name>] <previous-revision> <next-revision>
```

-h : Show usage information.

Example:

```
kubecompare deployment/my-app 3 4
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
```
