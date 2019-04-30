require 'pg'

class Model
  # pg: http://www.ownway.info/Ruby/pg/about
  def initialize()
    puts 'connect postgres'
    
    @conn = PG::connect(
        host:     'localhost',
        user:     'postgres',
        password: 'postgres',
        dbname:   'tag-mng'
    )
  end

  def select()
    q = 'SELECT name FROM users'
    rslt = @conn.exec(q)

    n = ''
    rslt.each do |row|
      n = row['name']
      puts n
    end
    n
  end
end