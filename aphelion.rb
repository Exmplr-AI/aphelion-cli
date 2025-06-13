class Aphelion < Formula
  desc "Command-line interface for Aphelion Gateway"
  homepage "https://github.com/Exmplr-AI/aphelion-cli"
  url "https://github.com/Exmplr-AI/aphelion-cli/archive/v#{version}.tar.gz"
  license "MIT"
  head "https://github.com/Exmplr-AI/aphelion-cli.git", branch: "main"

  depends_on "go" => :build

  def install
    ENV["CGO_ENABLED"] = "0"
    ldflags = %W[
      -s -w
      -X main.Version=#{version}
      -X main.GitCommit=#{tap.user}
      -X main.BuildDate=#{time.iso8601}
    ]

    system "go", "build", *std_go_args(ldflags: ldflags), "./main.go"

    # Install shell completions
    generate_completions_from_executable(bin/"aphelion", "completion")
  end

  test do
    assert_match "aphelion version #{version}", shell_output("#{bin}/aphelion version")
  end
end