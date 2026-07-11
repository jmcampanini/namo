class Namo < Formula
  desc "Generate memorable, sortable names"
  homepage "https://github.com/jmcampanini/namo"
  head "https://github.com/jmcampanini/namo.git", branch: "main"

  depends_on "go" => :build

  def install
    ldflags = %W[
      -s -w
      -X github.com/jmcampanini/namo/cmd.Version=#{version}
    ]
    system "go", "build", "-buildvcs=false", *std_go_args(ldflags:)
    generate_completions_from_executable(bin/"namo", "completion")
  end

  test do
    assert_match "namo version HEAD-", shell_output("#{bin}/namo --version")
  end
end
