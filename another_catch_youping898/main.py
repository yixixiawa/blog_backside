import requests
import time
from typing import List, Dict
from database.sqlite import ItemDatabase

class YouPin898Client:
    def __init__(self, authorization: str):
        self.session = requests.Session()
        self.session.headers.update({
            "Accept": "application/json, text/plain, */*",
            "Accept-Encoding": "gzip, deflate, br, zstd",
            "Accept-Language": "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7",
            "App-Version": "5.26.0",
            "AppVersion": "5.26.0",
            "Content-Type": "application/json",
            "Origin": "https://www.youpin898.com",
            "Referer": "https://www.youpin898.com/",
            "Sec-Fetch-Dest": "empty",
            "Sec-Fetch-Mode": "cors",
            "Sec-Fetch-Site": "same-site",
            "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36 Edg/142.0.0.0",
            "appType": "1",
            "platform": "pc",
            "secret-v": "h5_v1",
            "x-dev-access": "yes",
            "authorization": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJqdGkiOiIzN2UzZWNmOGMwMmU0ZGVlOTk4ZWI1NWYwODQwOWIxMyIsIm5hbWVpZCI6IjE4NjQyODMiLCJJZCI6IjE4NjQyODMiLCJ1bmlxdWVfbmFtZSI6IllQMDAwMTg2NDI4MyIsIk5hbWUiOiJZUDAwMDE4NjQyODMiLCJ2ZXJzaW9uIjoiZ2dIIiwibmJmIjoxNzYyNzQzMzA0LCJleHAiOjE3NjM2MDczMDQsImlzcyI6InlvdXBpbjg5OC5jb20iLCJkZXZpY2VJZCI6IjJmMjZlNGFiLWZlOGItNDcyZS1iOTQ5LTEzZWE5NWM4MzdhYyIsImF1ZCI6InVzZXIifQ.Y5KENnAoG0GkNlub5ThxF-lptvpIEDnNWlRZVT_BtYo",
            "deviceUk": "5FTCicCrcfkZqtqElju7XHiH80Q39udiiBvuby2J2lMq3haQNZbhsDJRryq5nwL1O",
            "uk": "5FRnK6XGgteNY6JsxDoBf7w4ZcZlNspS9IR0HfJujMKmNh2VaCQHaIep8Eybh0v1K"
        })
        self.url = "https://api.youpin898.com/api/homepage/pc/goods/market/querySaleTemplate"
    
    def fetch_page(self, page: int, page_size: int = 20) -> List[Dict]:
        """获取单页数据"""
        try:
            resp = self.session.post(
                self.url,
                json={"gameId": 730, "pageIndex": page, "pageSize": page_size},
                timeout=10
            )
            data = resp.json()
            return data.get("Data", []) if data.get("Code") == 0 else []
        except Exception as e:
            print(f"第 {page} 页失败: {e}")
            return []
    
    def fetch_and_save(self, start: int = 1, end: int = 1250, delay: float = 1.0):
        """批量爬取并保存到数据库"""
        db = ItemDatabase("./data/data.db")
        total = 0
        
        print(f"开始爬取第 {start}-{end} 页并保存到数据库...")
        
        try:
            for page in range(start, end + 1):
                items = self.fetch_page(page)
                if items:
                    count = db.insert_batch(items)
                    total += count
                    print(f"[{page}/{end}] 保存 {count} 条，累计 {total}")
                
                if page < end:
                    time.sleep(delay)
        except KeyboardInterrupt:
            print("\n用户中断，正在保存...")
        finally:
            db.close()
        
        print(f"完成！共保存 {total} 条数据")
    
    def fetch_and_update_smart(self, start: int = 1, end: int = 1250, delay: float = 1.0):
        """智能爬取并更新数据库（避免重复爬取静态数据）"""
        db = ItemDatabase("./data/data.db")
        total_new = 0
        total_updated = 0
        
        print(f"开始智能爬取第 {start}-{end} 页...")
        print("说明：新商品将完整保存，已存在商品只更新价格相关字段")
        
        try:
            for page in range(start, end + 1):
                items = self.fetch_page(page)
                if items:
                    result = db.upsert_batch_smart(items)
                    total_new += result['new']
                    total_updated += result['updated']
                    print(f"[{page}/{end}] 新增 {result['new']} 条，更新 {result['updated']} 条")
                
                if page < end:
                    time.sleep(delay)
        except KeyboardInterrupt:
            print("\n用户中断，正在保存...")
        finally:
            db.close()
        
        print(f"完成！新增 {total_new} 条，更新 {total_updated} 条数据")
    
    def fetch_all_with_auto_stop(self, start: int = 1, delay: float = 1.0, empty_pages_limit: int = 3, smart_update: bool = True):
        """自动爬取所有数据，遇到连续空页面时停止"""
        db = ItemDatabase("./data/data.db")
        total_new = 0
        total_updated = 0
        current_page = start
        consecutive_empty_pages = 0
        
        mode_desc = "智能更新模式" if smart_update else "完整替换模式"
        print(f"开始自动爬取，从第 {start} 页开始 ({mode_desc})")
        print(f"连续 {empty_pages_limit} 页无数据时将自动停止")
        
        if smart_update:
            print("智能模式：新商品完整保存，已存在商品只更新价格字段")
        
        try:
            while consecutive_empty_pages < empty_pages_limit:
                print(f"正在爬取第 {current_page} 页...")
                items = self.fetch_page(current_page)
                
                if items:
                    # 有数据，重置空页面计数器
                    consecutive_empty_pages = 0
                    
                    if smart_update:
                        result = db.upsert_batch_smart(items)
                        total_new += result['new']
                        total_updated += result['updated']
                        print(f"[第 {current_page} 页] 新增 {result['new']} 条，更新 {result['updated']} 条")
                    else:
                        count = db.insert_batch(items)
                        total_new += count
                        print(f"[第 {current_page} 页] 保存 {count} 条数据")
                else:
                    # 无数据，增加空页面计数器
                    consecutive_empty_pages += 1
                    print(f"[第 {current_page} 页] 无数据，连续空页面: {consecutive_empty_pages}/{empty_pages_limit}")
                
                current_page += 1
                
                # 如果还需要继续且未达到空页面限制，则等待
                if consecutive_empty_pages < empty_pages_limit:
                    time.sleep(delay)
                    
        except KeyboardInterrupt:
            print("\n用户中断，正在保存...")
        finally:
            db.close()
        
        if consecutive_empty_pages >= empty_pages_limit:
            print(f"检测到连续 {empty_pages_limit} 页无数据，自动停止")
        
        if smart_update:
            print(f"爬取完成！共爬取 {current_page - start} 页，新增 {total_new} 条，更新 {total_updated} 条数据")
        else:
            print(f"爬取完成！共爬取 {current_page - start} 页，保存 {total_new} 条数据")

# 使用示例
if __name__ == "__main__":
    TOKEN = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9..."  # 请替换为实际的token
    
    client = YouPin898Client(TOKEN)
    
    print("=== 智能更新模式 ===")
    client.fetch_all_with_auto_stop(start=1, delay=0.5, empty_pages_limit=3, smart_update=True)
    print("所有操作完成！")