class Ctp < Formula
    desc "[CPT] Code to Prompt - A tool to generate a prompt from codes"
    homepage "https://github.com/alesr/codetoprompt"
    url "https://github.com/alesr/codetoprompt/archive/refs/tags/v1.0.0.tar.gz"
    sha256 "b974a10de075e9e9c950a6ed3d541704136bed5b4d42396cfb6c095907a1614f"
    license "MIT"
    head "https://github.com/alesr/codetoprompt.git", branch: "master"
  
    bottle do
      sha256 cellar: :any_skip_relocation, arm64_ventura:  "a713440029965885a313b22d7fba78b30b2e56003a2b2955f8dfc01029e8836a"
      sha256 cellar: :any_skip_relocation, arm64_monterey: "a713440029965885a313b22d7fba78b30b2e56003a2b2955f8dfc01029e8836a"
      sha256 cellar: :any_skip_relocation, arm64_big_sur:  "a713440029965885a313b22d7fba78b30b2e56003a2b2955f8dfc01029e8836a"
      sha256 cellar: :any_skip_relocation, ventura:        "ba3883ee8187e4990fba2df1315831f211e579ecd83f680f582c9f33af541a34"
      sha256 cellar: :any_skip_relocation, monterey:       "ba3883ee8187e4990fba2df1315831f211e579ecd83f680f582c9f33af541a34"
      sha256 cellar: :any_skip_relocation, big_sur:        "ba3883ee8187e4990fba2df1315831f211e579ecd83f680f582c9f33af541a34"
      sha256 cellar: :any_skip_relocation, catalina:       "ba3883ee8187e4990fba2df1315831f211e579ecd83f680f582c9f33af541a34"
      sha256 cellar: :any_skip_relocation, x86_64_linux:   "a713440029965885a313b22d7fba78b30b2e56003a2b2955f8dfc01029e8836a"
    end

    depends_on "go" => :build
  
    def install
      ENV["GOPATH"] = buildpath
      # Move the contents of the repo to another directory to avoid `go mod` errors
      (buildpath/"src/github.com/alesr/codetoprompt").install buildpath.children
      cd "src/github.com/alesr/codetoprompt" do
        system "go", "mod", "download"
        system "go", "build", "-o", bin/"ctp"
        prefix.install_metafiles
      end
    end
  
    test do
      system "#{bin}/ctp", "--help"
    end
  end
