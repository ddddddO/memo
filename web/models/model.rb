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

  def list(user_id)
    q = 'SELECT id, subject FROM memos WHERE users_id=$1'
    rslt = @conn.exec(q, [user_id])

    rows = []
    rslt.each do |row|
      rows.push(row)
    end
    rows
  end

  def detail(memo_id, user_id)
    q = 'SELECT id, subject, content FROM memos WHERE id=$1 AND users_id=$2'
    rslt = @conn.exec(q, [memo_id, user_id])

    rows = []
    rslt.each do |row|
      rows.push(row)
    end
    rows
  end

  def update(args)
    # 一時的
    if args['memo_id'].empty?
      args['memo_id'] = "3"
    end

    q = <<~EOS
      INSERT INTO memos(id, subject, content, users_id)
      VALUES($1, $2, $3, $4)
      ON CONFLICT(id)
      DO UPDATE SET subject=$5, content=$6
      RETURNING id
    EOS

    rslt = @conn.exec(q,[
        args['memo_id'],
        args['subject'],
        args['content'],
        args['user_id'],
        args['subject'],
        args['content']
    ])
    
    return rslt[0]['id']
  end
end