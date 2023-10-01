class Kubecompare < Formula
    desc "A tool to compare Kubernetes rollout histories"
    url "https://github.com/wozniakpl/kubecompare/archive/v0.1.0.tar.gz"
    sha256 "9bfd5bc94deebb9c814be9f40b5bd928a73e1ba0535a57674ea249d2f7e1a619"
    license "MIT"
  
    depends_on "go" => :build
  
    def install
      system "go", "build", *std_go_args, "-o", bin/"kubecompare"
    end
  
    test do
      system "#{bin}/kubecompare", "-h"
    end
  end
  