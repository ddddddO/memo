require 'sinatra'
require 'sinatra/base'
require 'sinatra/reloader'

get '/' do
  'sinatra!!!!!'
end

get '/tmp_erb' do
  erb :tmp_erb
end
