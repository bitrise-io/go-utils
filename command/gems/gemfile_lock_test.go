package gems

import (
	"reflect"
	"testing"
)

func Test_ParseFastlaneVersion(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name               string
		gemfileLockContent string
		wantGemVersion     GemVersion
		wantErr            bool
	}{
		{
			gemfileLockContent: gemfileLockContent,
			wantGemVersion: GemVersion{
				Version: "2.13.0",
				Found:   true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGemVersion, err := ParseFastlaneVersion(tt.gemfileLockContent)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFastlaneVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotGemVersion, tt.wantGemVersion) {
				t.Errorf("ParseFastlaneVersion() = %v, want %v", gotGemVersion, tt.wantGemVersion)
			}
		})
	}
}

func Test_ParseBundlerVersion(t *testing.T) {
	tests := []struct {
		name               string
		gemfileLockContent string
		want               GemVersion
		wantErr            bool
	}{
		{
			name:               "should match",
			gemfileLockContent: gemfileLockContent,
			want: GemVersion{
				Version: "1.13.6",
				Found:   true,
			},
		},
		{
			name: "newline after version",
			gemfileLockContent: `BUNDLED WITH
      1.13.6
      
      `,
			want: GemVersion{
				Version: "1.13.6",
				Found:   true,
			},
		},
		{
			name: "newline before version",
			gemfileLockContent: `BUNDLED WITH
      
      1.13.6`,
			want: GemVersion{
				Version: "1.13.6",
				Found:   true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseBundlerVersion(tt.gemfileLockContent)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseBundlerVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseBundlerVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

const gemfileLockContent = `GIT
  remote: git://xyz.git
  revision: xyz
  branch: patch-1
  specs:
    fastlane-xyz (1.0.2)

GEM
  remote: https://rubygems.org/
  specs:
    CFPropertyList (2.3.5)
    activesupport (4.2.7.1)
      i18n (~> 0.7)
      json (~> 1.7, >= 1.7.7)
      minitest (~> 5.1)
      thread_safe (~> 0.3, >= 0.3.4)
      tzinfo (~> 1.1)
    addressable (2.5.0)
      public_suffix (~> 2.0, >= 2.0.2)
    babosa (1.0.2)
    claide (1.0.1)
    cocoapods (1.1.1)
      activesupport (>= 4.0.2, < 5)
      claide (>= 1.0.1, < 2.0)
      cocoapods-core (= 1.1.1)
      cocoapods-deintegrate (>= 1.0.1, < 2.0)
      cocoapods-downloader (>= 1.1.2, < 2.0)
      cocoapods-plugins (>= 1.0.0, < 2.0)
      cocoapods-search (>= 1.0.0, < 2.0)
      cocoapods-stats (>= 1.0.0, < 2.0)
      cocoapods-trunk (>= 1.1.1, < 2.0)
      cocoapods-try (>= 1.1.0, < 2.0)
      colored (~> 1.2)
      escape (~> 0.0.4)
      fourflusher (~> 2.0.1)
      gh_inspector (~> 1.0)
      molinillo (~> 0.5.1)
      nap (~> 1.0)
      xcodeproj (>= 1.3.3, < 2.0)
    cocoapods-core (1.1.1)
      activesupport (>= 4.0.2, < 5)
      fuzzy_match (~> 2.0.4)
      nap (~> 1.0)
    cocoapods-deintegrate (1.0.1)
    cocoapods-downloader (1.1.3)
    cocoapods-plugins (1.0.0)
      nap
    cocoapods-search (1.0.0)
    cocoapods-stats (1.0.0)
    cocoapods-trunk (1.1.2)
      nap (>= 0.8, < 2.0)
      netrc (= 0.7.8)
    cocoapods-try (1.1.0)
    colored (1.2)
    commander (4.4.3)
      highline (~> 1.7.2)
    domain_name (0.5.20161129)
      unf (>= 0.0.5, < 1.0.0)
    dotenv (2.2.0)
    escape (0.0.4)
    excon (0.54.0)
    faraday (0.11.0)
      multipart-post (>= 1.2, < 3)
    faraday-cookie_jar (0.0.6)
      faraday (>= 0.7.4)
      http-cookie (~> 1.0.0)
    faraday_middleware (0.11.0.1)
      faraday (>= 0.7.4, < 1.0)
    fastimage (2.0.1)
      addressable (~> 2)
    fastlane (2.13.0)
      activesupport (< 5)
      addressable (>= 2.3, < 3.0.0)
      babosa (>= 1.0.2, < 2.0.0)
      bundler (>= 1.12.0, < 2.0.0)
      colored
      commander (>= 4.4.0, < 5.0.0)
      dotenv (>= 2.1.1, < 3.0.0)
      excon (>= 0.45.0, < 1.0.0)
      faraday (~> 0.9)
      faraday-cookie_jar (~> 0.0.6)
      faraday_middleware (~> 0.9)
      fastimage (>= 1.6)
      gh_inspector (>= 1.0.1, < 2.0.0)
      google-api-client (~> 0.9.2)
      highline (>= 1.7.2, < 2.0.0)
      json (< 3.0.0)
      mini_magick (~> 4.5.1)
      multi_json
      multi_xml (~> 0.5)
      multipart-post (~> 2.0.0)
      plist (>= 3.1.0, < 4.0.0)
      rubyzip (>= 1.1.0, < 2.0.0)
      security (= 0.1.3)
      slack-notifier (>= 1.3, < 2.0.0)
      terminal-notifier (>= 1.6.2, < 2.0.0)
      terminal-table (>= 1.4.5, < 2.0.0)
      word_wrap (~> 1.0.0)
      xcodeproj (>= 0.20, < 2.0.0)
      xcpretty (>= 0.2.4, < 1.0.0)
      xcpretty-travis-formatter (>= 0.0.3)
    fastlane-plugin-tpa (1.1.0)
    fourflusher (2.0.1)
    fuzzy_match (2.0.4)
    gh_inspector (1.0.3)
    google-api-client (0.9.26)
      addressable (~> 2.3)
      googleauth (~> 0.5)
      httpclient (~> 2.7)
      hurley (~> 0.1)
      memoist (~> 0.11)
      mime-types (>= 1.6)
      representable (~> 2.3.0)
      retriable (~> 2.0)
    googleauth (0.5.1)
      faraday (~> 0.9)
      jwt (~> 1.4)
      logging (~> 2.0)
      memoist (~> 0.12)
      multi_json (~> 1.11)
      os (~> 0.9)
      signet (~> 0.7)
    highline (1.7.8)
    http-cookie (1.0.3)
      domain_name (~> 0.5)
    httpclient (2.8.3)
    hurley (0.2)
    i18n (0.8.0)
    json (1.8.6)
    jwt (1.5.6)
    little-plugger (1.1.4)
    logging (2.1.0)
      little-plugger (~> 1.1)
      multi_json (~> 1.10)
    memoist (0.15.0)
    mime-types (3.1)
      mime-types-data (~> 3.2015)
    mime-types-data (3.2016.0521)
    mini_magick (4.5.1)
    minitest (5.10.1)
    molinillo (0.5.5)
    multi_json (1.12.1)
    multi_xml (0.6.0)
    multipart-post (2.0.0)
    nanaimo (0.2.3)
    nap (1.1.0)
    netrc (0.7.8)
    os (0.9.6)
    plist (3.2.0)
    public_suffix (2.0.5)
    representable (2.3.0)
      uber (~> 0.0.7)
    retriable (2.1.0)
    rouge (1.11.1)
    rubyzip (1.2.0)
    security (0.1.3)
    signet (0.7.3)
      addressable (~> 2.3)
      faraday (~> 0.9)
      jwt (~> 1.5)
      multi_json (~> 1.10)
    slack-notifier (1.5.1)
    terminal-notifier (1.7.1)
    terminal-table (1.7.3)
      unicode-display_width (~> 1.1.1)
    thread_safe (0.3.5)
    tzinfo (1.2.2)
      thread_safe (~> 0.1)
    uber (0.0.15)
    unf (0.1.4)
      unf_ext
    unf_ext (0.0.7.2)
    unicode-display_width (1.1.3)
    word_wrap (1.0.0)
    xcode-install (2.1.1)
      claide (>= 0.9.1, < 1.1.0)
      fastlane (>= 2.1.0, < 3.0.0)
    xcodeproj (1.4.2)
      CFPropertyList (~> 2.3.3)
      activesupport (>= 3)
      claide (>= 1.0.1, < 2.0)
      colored (~> 1.2)
      nanaimo (~> 0.2.3)
    xcpretty (0.2.4)
> 1.8)
    xcpretty-travis-formatter (0.0.4)
      xcpretty (~> 0.2, >= 0.0.7)

PLATFORMS
  ruby

DEPENDENCIES
  cocoapods (~> 1.1.0)
  dotenv (~> 2.0)
  fastlane (~> 2.0)
  fastlane-plugin-prepare_build_resources!
  fastlane-plugin-tpa (~> 1.1.0)
  xcode-install

BUNDLED WITH
   1.13.6`
