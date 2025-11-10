import socket
import json
import threading
from typing import Dict, List
from database.sqlite import ItemDatabase

class DatabaseServer:
    """数据库 Socket 服务器"""
    
    def __init__(self, host: str = '0.0.0.0', port: int = 8080, db_path: str = "./data/data.db"):
        self.host = host
        self.port = port
        self.db_path = db_path
        self.socket = None
        self.running = False
    
    def start(self):
        """启动服务器"""
        try:
            # 创建 Socket
            self.socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            self.socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
            
            # 绑定地址
            self.socket.bind((self.host, self.port))
            self.socket.listen(5)
            self.running = True
            
            print(f"🚀 数据库服务器启动成功")
            print(f"   地址: {self.host}:{self.port}")
            print(f"   数据库: {self.db_path}")
            print(f"   等待客户端连接...")
            
            while self.running:
                try:
                    # 接受连接
                    client_socket, client_addr = self.socket.accept()
                    print(f"📱 客户端连接: {client_addr}")
                    
                    # 创建线程处理客户端
                    thread = threading.Thread(
                        target=self.handle_client,
                        args=(client_socket, client_addr)
                    )
                    thread.daemon = True
                    thread.start()
                    
                except socket.error as e:
                    if self.running:
                        print(f"❌ Socket 错误: {e}")
                    break
                    
        except Exception as e:
            print(f"❌ 服务器启动失败: {e}")
        finally:
            self.stop()
    
    def handle_client(self, client_socket: socket.socket, client_addr):
        """处理客户端请求"""
        db = None
        try:
            # 连接数据库
            db = ItemDatabase(self.db_path)
            
            while True:
                # 接收请求
                data = client_socket.recv(4096).decode('utf-8')
                if not data:
                    break
                
                print(f"📨 收到请求: {client_addr} -> {data[:100]}...")
                
                # 处理请求
                try:
                    request = json.loads(data)
                    response = self.process_request(db, request)
                except json.JSONDecodeError:
                    response = {"error": "无效的 JSON 格式"}
                except Exception as e:
                    response = {"error": str(e)}
                
                # 发送响应
                response_json = json.dumps(response, ensure_ascii=False)
                client_socket.send(response_json.encode('utf-8'))
                
        except Exception as e:
            print(f"❌ 处理客户端 {client_addr} 出错: {e}")
        finally:
            if db:
                db.close()
            client_socket.close()
            print(f"🔌 客户端断开: {client_addr}")
    
    def process_request(self, db: ItemDatabase, request: Dict) -> Dict:
        """处理数据库请求"""
        action = request.get('action')
        params = request.get('params', {})
        
        try:
            if action == 'stats':
                # 获取统计信息
                return {"success": True, "data": db.get_stats()}
            
            elif action == 'cheapest':
                # 获取最便宜的饰品
                limit = params.get('limit', 10)
                return {"success": True, "data": db.get_cheapest(limit)}
            
            elif action == 'search_name':
                # 按名称搜索
                name = params.get('name', '')
                if not name:
                    return {"success": False, "error": "缺少 name 参数"}
                return {"success": True, "data": db.query_by_name(name)}
            
            elif action == 'price_range':
                # 按价格区间查询
                min_price = params.get('min_price', 0)
                max_price = params.get('max_price', 999999)
                return {"success": True, "data": db.query_by_price_range(min_price, max_price)}
            
            elif action == 'all':
                # 获取所有数据（谨慎使用）
                limit = params.get('limit', 100)
                db.cursor.execute("SELECT * FROM items LIMIT ?", (limit,))
                results = [dict(row) for row in db.cursor.fetchall()]
                return {"success": True, "data": results}
            
            else:
                return {"success": False, "error": f"未知操作: {action}"}
                
        except Exception as e:
            return {"success": False, "error": str(e)}
    
    def stop(self):
        """停止服务器"""
        self.running = False
        if self.socket:
            self.socket.close()
        print("🛑 服务器已停止")

# 客户端测试函数
def test_client(host: str = 'localhost', port: int = 8080):
    """测试客户端"""
    try:
        # 连接服务器
        client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        client.connect((host, port))
        print(f"✅ 连接到服务器: {host}:{port}")
        
        # 测试请求
        test_requests = [
            {"action": "stats"},
            {"action": "cheapest", "params": {"limit": 5}},
            {"action": "search_name", "params": {"name": "AK"}},
            {"action": "price_range", "params": {"min_price": 0.1, "max_price": 1.0}}
        ]
        
        for i, req in enumerate(test_requests, 1):
            print(f"\n📤 测试 {i}: {req['action']}")
            
            # 发送请求
            client.send(json.dumps(req).encode('utf-8'))
            
            # 接收响应
            response = client.recv(4096).decode('utf-8')
            data = json.loads(response)
            
            if data.get('success'):
                results = data['data']
                if isinstance(results, list):
                    print(f"✅ 成功，返回 {len(results)} 条记录")
                    for item in results[:3]:  # 显示前3条
                        if 'commodity_name' in item:
                            print(f"   - {item['commodity_name']}: ¥{item.get('price', 'N/A')}")
                        else:
                            print(f"   - {item}")
                else:
                    print(f"✅ 成功: {results}")
            else:
                print(f"❌ 失败: {data.get('error')}")
        
        client.close()
        print("\n🔌 客户端断开连接")
        
    except Exception as e:
        print(f"❌ 客户端错误: {e}")

# 启动函数
def create_sock(host: str = '0.0.0.0', port: int = 8080):
    """创建并启动数据库服务器"""
    server = DatabaseServer(host, port)
    
    try:
        server.start()
    except KeyboardInterrupt:
        print("\n⚠️  收到中断信号...")
    finally:
        server.stop()

if __name__ == "__main__":
    import sys
    
    if len(sys.argv) > 1 and sys.argv[1] == 'test':
        # 测试模式
        print("🧪 启动客户端测试...")
        test_client()
    else:
        # 服务器模式
        print("🚀 启动数据库服务器...")
        create_sock()