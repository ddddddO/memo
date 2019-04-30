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

  def login(name, passwd)
    q = 'SELECT id FROM users WHERE name=$1 AND passwd=$2'
    rslt = @conn.exec(q, [name, passwd])
    
    # TODO: 要リファクタ&認証エラーハンドル(401)
    if !rslt.nil?
      user_id = ''
      rslt.each do |row|
        user_id = row['id']
      end

      if !user_id.empty?
        return user_id
      end
    else
      return 'failed to login'
    end
    return 'failed to login'
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