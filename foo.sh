#!/usr/bin/env ruby
n = rand(20)
beginning_time = Time.now
if ARGV[0] == 13
  puts 'Sorry buddy'
  exit 127
end
puts "Staring event job ##{ARGV[0]}"
sleep ARGV[0].to_i*2
puts "Event job ##{ARGV[0]} done after #{ Time.now - beginning_time} ms"
