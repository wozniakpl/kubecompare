# kubecompare

A CLI tool to compare different revisions of Kubernetes resources.

## Description

Have you ever updated a Kubernetes deployment and wondered what changed between the previous and current version? `kubecompare` is a command-line utility that helps users compare different revisions of Kubernetes deployments, daemonsets, or statefulsets to understand what changed between versions.

## Installation

### Homebrew

If you have Homebrew installed, you can install kubecompare from our custom tap:

```
brew tap wozniakpl/kubecompare
brew install kubecompare
```

### From Source

Clone the repository and build the application using Go:

```bash
git clone https://github.com/wozniakpl/kubecompare.git
cd kubecompare
go build -o kubecompare
```

## Usage

```
kubecompare [ <resource-type> <resource-name> | <resource-type>/<resource-name> ] [ <previous-revision> <next-revision> ]
```

### Example:

First, let's create some resource. Let's say it's a deployment:

```bash
$ kubectl create deployment nginx --image=nginx
```

Now, let's change the image in the deployment:

```bash
$ kubectl edit deployment nginx
```

Once it's done, we can see that we have two revisions available:

```
$ kubecompare deployment nginx
deployment.apps/nginx
REVISION  CHANGE-CAUSE
1         <none>
2         <none>
```

And now we can see the diff:

```
$ kubecompare deployment nginx 1 2
--- /var/folders/ck/1k1429w513v98xw2r27pqmzw0000gn/T/revision3918300321 2023-10-02 21:04:46
+++ /var/folders/ck/1k1429w513v98xw2r27pqmzw0000gn/T/revision167025667  2023-10-02 21:04:46
@@ -1,10 +1,10 @@
-deployment.apps/nginx with revision #1
+deployment.apps/nginx with revision #2
 Pod Template:
   Labels:      app=nginx
-       pod-template-hash=748c667d99
+       pod-template-hash=5f6f57cd9
   Containers:
    nginx:
-    Image:     nginx
+    Image:     busybox
     Port:      <none>
     Host Port: <none>
     Environment:       <none>
```


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
