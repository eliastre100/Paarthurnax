#!/usr/bin/env ruby
# frozen_string_literal: true

ENV["BUNDLE_GEMFILE"] ||= File.join(__dir__, "Gemfile")
require "bundler/setup"
require "i18n"
require "json"
require "rails-i18n"

I18n.enforce_available_locales = false

gem_root = Gem::Specification.find_by_name("rails-i18n").full_gem_path
Dir[File.join(gem_root, "lib", "rails_i18n", "common_pluralizations", "*.rb")].sort.each do |file|
  require file
end
require File.join(gem_root, "lib", "rails_i18n", "pluralization.rb")

I18n.load_path = Dir[File.join(gem_root, "rails", "pluralization", "*.rb")]
I18n.backend.load_translations

PLURALIZATION_RULES = {
    # Common pluralization rules
    east_slavic: ::RailsI18n::Pluralization::EastSlavic.rule,
    one_few_other: ::RailsI18n::Pluralization::OneFewOther.rule,
    one_other: ::RailsI18n::Pluralization::OneOther.rule,
    one_two_other: ::RailsI18n::Pluralization::OneTwoOther.rule,
    one_upto_two_other: ::RailsI18n::Pluralization::OneUptoTwoOther.rule,
    one_with_zero_other: ::RailsI18n::Pluralization::OneWithZeroOther.rule,
    other: ::RailsI18n::Pluralization::Other.rule,
    romanian: ::RailsI18n::Pluralization::Romanian.rule,
    west_slavic: ::RailsI18n::Pluralization::WestSlavic.rule,

    # pluralization.rb
    arabic: ::RailsI18n::Pluralization::Arabic.rule,
    scottish_gaelic: ::RailsI18n::Pluralization::ScottishGaelic.rule,
    upper_sorbian: ::RailsI18n::Pluralization::UpperSorbian.rule,
    lithuanian: ::RailsI18n::Pluralization::Lithuanian.rule,
    latvian: ::RailsI18n::Pluralization::Latvian.rule,
    macedonian: ::RailsI18n::Pluralization::Macedonian.rule,
    polish: ::RailsI18n::Pluralization::Polish.rule,
    Slovenian: ::RailsI18n::Pluralization::Slovenian.rule
}

PLURALIZATION_RULES_MAP = PLURALIZATION_RULES.map do |key, rule|
    [rule.source_location, key]
end.to_h


rules = I18n.backend.send(:translations).filter_map do |locale, translations|
  pluralization = translations.dig(:i18n, :plural)
  next unless pluralization

  rule_name = PLURALIZATION_RULES_MAP[pluralization.fetch(:rule).source_location]
  [locale, { keys: pluralization.fetch(:keys).map(&:to_s), rule_name: rule_name }]
end.to_h

puts JSON.generate(rules)
