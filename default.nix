{buildGoModule}:
buildGoModule {
  pname = "generator";
  version = "0.0.1";

  src = ./.;

  vendorHash = "sha256-VstMR+lm1C4pr2VaM3S35xrMu2TvTDuqzS+fHi7Kz6k=";

  ldflags = ["-s" "-w"];
}
