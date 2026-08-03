#!/usr/bin/env ruby
# frozen_string_literal: true

ENV["BUNDLE_GEMFILE"] ||= File.join(__dir__, "Gemfile")
require "bundler/setup"
require "i18n"
require "json"
require "rails-i18n"

requests = JSON.parse(ARGF.read)

I18n.enforce_available_locales = false

gem_root = Gem::Specification.find_by_name("rails-i18n").full_gem_path
I18n.load_path = Dir[File.join(gem_root, "rails", "pluralization", "*.rb")]
I18n.backend.load_translations

keys = requests.map do |request|
  locale = request.fetch("locale")
  hint = request.fetch("hint").to_i
  pluralization = I18n.backend.send(:translations).dig(locale.to_sym, :i18n, :plural)
  abort "no pluralization rule for #{locale.inspect}" unless pluralization

  pluralization.fetch(:rule).call(hint).to_s
end

puts JSON.generate(keys)
