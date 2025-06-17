# This is a template Homebrew formula for aphelion-cli
# This file should be placed in a homebrew tap repository: homebrew-aphelion/Formula/aphelion.rb

class Aphelion < Formula
  desc "Command-line interface for the Aphelion Gateway platform"
  homepage "https://github.com/Exmplr-AI/aphelion-cli"
  version "1.0.0" # This will be updated by the release workflow
  
  if OS.mac?
    if Hardware::CPU.arm?
      url "https://github.com/Exmplr-AI/aphelion-cli/releases/download/v#{version}/aphelion-darwin-arm64"
      sha256 "REPLACE_WITH_ARM64_SHA256" # This will be updated by release workflow
    else
      url "https://github.com/Exmplr-AI/aphelion-cli/releases/download/v#{version}/aphelion-darwin-amd64"
      sha256 "REPLACE_WITH_AMD64_SHA256" # This will be updated by release workflow
    end
  elsif OS.linux?
    if Hardware::CPU.arm?
      url "https://github.com/Exmplr-AI/aphelion-cli/releases/download/v#{version}/aphelion-linux-arm64"
      sha256 "REPLACE_WITH_LINUX_ARM64_SHA256" # This will be updated by release workflow
    else
      url "https://github.com/Exmplr-AI/aphelion-cli/releases/download/v#{version}/aphelion-linux-amd64"
      sha256 "REPLACE_WITH_LINUX_AMD64_SHA256" # This will be updated by release workflow
    end
  end

  def install
    bin.install Dir["*"].first => "aphelion"
    
    # Generate and install shell completions
    output = Utils.safe_popen_read(bin/"aphelion", "completion", "bash")
    (bash_completion/"aphelion").write output
    
    output = Utils.safe_popen_read(bin/"aphelion", "completion", "zsh")
    (zsh_completion/"_aphelion").write output
    
    output = Utils.safe_popen_read(bin/"aphelion", "completion", "fish")
    (fish_completion/"aphelion.fish").write output
  end

  test do
    system "#{bin}/aphelion", "version"
    assert_match version.to_s, shell_output("#{bin}/aphelion version")
  end
end