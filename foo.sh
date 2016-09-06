#!/usr/bin/env ruby
n = rand(20)
if ARGV[0] == 13
  puts 'Sorry buddy'
  exit 127
end
sleep n
puts "Job ##{ARGV[0]} done after #{n}"
