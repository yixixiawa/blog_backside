import sqlite3
import os
from typing import List, Dict, Optional
from datetime import datetime


class ItemDatabase:
    """CS:GO 饰品数据库管理类"""
    
    def __init__(self, db_path: str = "./data/data.db"):
        """初始化数据库连接"""
        # 确保 data 目录存在
        os.makedirs(os.path.dirname(db_path), exist_ok=True)
        
        self.db_path = db_path
        self.conn = sqlite3.connect(db_path)
        self.conn.row_factory = sqlite3.Row
        self.cursor = self.conn.cursor()
        
        # 创建表
        self._create_tables()
        print(f"数据库已连接: {db_path}")
    
    def _create_tables(self):
        """创建数据表"""
        # 商品表
        self.cursor.execute("""
            CREATE TABLE IF NOT EXISTS items (
                id INTEGER PRIMARY KEY,
                game_id INTEGER,
                game_name TEXT,
                commodity_name TEXT NOT NULL,
                commodity_hash_name TEXT,
                icon_url TEXT,
                on_sale_count INTEGER,
                price REAL,
                steam_price REAL,
                steam_usd_price REAL,
                type_name TEXT,
                exterior TEXT,
                exterior_color TEXT,
                rarity TEXT,
                rarity_color TEXT,
                quality TEXT,
                quality_color TEXT,
                have_lease INTEGER,
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
            )
        """)
        
        # 创建索引
        self.cursor.execute("""
            CREATE INDEX IF NOT EXISTS idx_commodity_name 
            ON items(commodity_name)
        """)
        self.cursor.execute("""
            CREATE INDEX IF NOT EXISTS idx_price 
            ON items(price)
        """)
        self.cursor.execute("""
            CREATE INDEX IF NOT EXISTS idx_updated 
            ON items(updated_at)
        """)
        
        self.conn.commit()
        print("数据表已创建")
    
    def insert_item(self, item: Dict) -> bool:
        """插入单条数据"""
        try:
            self.cursor.execute("""
                INSERT OR REPLACE INTO items (
                    id, game_id, game_name, commodity_name, commodity_hash_name,
                    icon_url, on_sale_count, price, steam_price, steam_usd_price,
                    type_name, exterior, exterior_color, rarity, rarity_color,
                    quality, quality_color, have_lease, updated_at
                ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
            """, (
                item.get('id'),
                item.get('gameId'),
                item.get('gameName'),
                item.get('commodityName'),
                item.get('commodityHashName'),
                item.get('iconUrl'),
                item.get('onSaleCount'),
                self._to_float(item.get('price')),
                self._to_float(item.get('steamPrice')),
                self._to_float(item.get('steamUsdPrice')),
                item.get('typeName'),
                item.get('exterior'),
                item.get('exteriorColor'),
                item.get('rarity'),
                item.get('rarityColor'),
                item.get('quality'),
                item.get('qualityColor'),
                item.get('haveLease'),
                datetime.now()
            ))
            return True
        except Exception as e:
            print(f"插入失败: {e}")
            return False
    
    def insert_batch(self, items: List[Dict]) -> int:
        """批量插入数据"""
        success_count = 0
        for item in items:
            if self.insert_item(item):
                success_count += 1
        
        self.conn.commit()
        print(f"批量插入完成: {success_count}/{len(items)}")
        return success_count
    
    def upsert_item_smart(self, item: Dict) -> bool:
        """智能插入或更新商品数据
        - 如果商品不存在，插入完整数据
        - 如果商品已存在，只更新价格相关字段
        """
        try:
            item_id = item.get('id')
            if not item_id:
                return False
            
            # 检查商品是否已存在
            self.cursor.execute("SELECT id FROM items WHERE id = ?", (item_id,))
            exists = self.cursor.fetchone()
            
            if exists:
                # 商品已存在，只更新价格相关字段
                self.cursor.execute("""
                    UPDATE items SET
                        on_sale_count = ?,
                        price = ?,
                        steam_price = ?,
                        steam_usd_price = ?,
                        updated_at = ?
                    WHERE id = ?
                """, (
                    item.get('onSaleCount'),
                    self._to_float(item.get('price')),
                    self._to_float(item.get('steamPrice')),
                    self._to_float(item.get('steamUsdPrice')),
                    datetime.now(),
                    item_id
                ))
            else:
                # 商品不存在，插入完整数据
                self.cursor.execute("""
                    INSERT INTO items (
                        id, game_id, game_name, commodity_name, commodity_hash_name,
                        icon_url, on_sale_count, price, steam_price, steam_usd_price,
                        type_name, exterior, exterior_color, rarity, rarity_color,
                        quality, quality_color, have_lease, updated_at
                    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
                """, (
                    item.get('id'),
                    item.get('gameId'),
                    item.get('gameName'),
                    item.get('commodityName'),
                    item.get('commodityHashName'),
                    item.get('iconUrl'),
                    item.get('onSaleCount'),
                    self._to_float(item.get('price')),
                    self._to_float(item.get('steamPrice')),
                    self._to_float(item.get('steamUsdPrice')),
                    item.get('typeName'),
                    item.get('exterior'),
                    item.get('exteriorColor'),
                    item.get('rarity'),
                    item.get('rarityColor'),
                    item.get('quality'),
                    item.get('qualityColor'),
                    item.get('haveLease'),
                    datetime.now()
                ))
            
            return True
        except Exception as e:
            print(f"智能插入/更新失败: {e}")
            return False
    
    def upsert_batch_smart(self, items: List[Dict]) -> Dict[str, int]:
        """批量智能插入或更新数据"""
        new_count = 0
        updated_count = 0
        failed_count = 0
        
        for item in items:
            item_id = item.get('id')
            if not item_id:
                failed_count += 1
                continue
                
            try:
                # 检查商品是否已存在
                self.cursor.execute("SELECT id FROM items WHERE id = ?", (item_id,))
                exists = self.cursor.fetchone()
                
                if exists:
                    # 更新现有商品
                    self.cursor.execute("""
                        UPDATE items SET
                            on_sale_count = ?,
                            price = ?,
                            steam_price = ?,
                            steam_usd_price = ?,
                            updated_at = ?
                        WHERE id = ?
                    """, (
                        item.get('onSaleCount'),
                        self._to_float(item.get('price')),
                        self._to_float(item.get('steamPrice')),
                        self._to_float(item.get('steamUsdPrice')),
                        datetime.now(),
                        item_id
                    ))
                    updated_count += 1
                else:
                    # 插入新商品
                    self.cursor.execute("""
                        INSERT INTO items (
                            id, game_id, game_name, commodity_name, commodity_hash_name,
                            icon_url, on_sale_count, price, steam_price, steam_usd_price,
                            type_name, exterior, exterior_color, rarity, rarity_color,
                            quality, quality_color, have_lease, updated_at
                        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
                    """, (
                        item.get('id'),
                        item.get('gameId'),
                        item.get('gameName'),
                        item.get('commodityName'),
                        item.get('commodityHashName'),
                        item.get('iconUrl'),
                        item.get('onSaleCount'),
                        self._to_float(item.get('price')),
                        self._to_float(item.get('steamPrice')),
                        self._to_float(item.get('steamUsdPrice')),
                        item.get('typeName'),
                        item.get('exterior'),
                        item.get('exteriorColor'),
                        item.get('rarity'),
                        item.get('rarityColor'),
                        item.get('quality'),
                        item.get('qualityColor'),
                        item.get('haveLease'),
                        datetime.now()
                    ))
                    new_count += 1
                    
            except Exception as e:
                print(f"处理商品 {item_id} 失败: {e}")
                failed_count += 1
        
        self.conn.commit()
        
        result = {
            'new': new_count,
            'updated': updated_count,
            'failed': failed_count,
            'total': new_count + updated_count
        }
        
        print(f"批量处理完成: 新增 {new_count} 条，更新 {updated_count} 条，失败 {failed_count} 条")
        return result
    
    def query_by_name(self, name: str) -> List[Dict]:
        """按名称查询"""
        self.cursor.execute("""
            SELECT * FROM items 
            WHERE commodity_name LIKE ? 
            ORDER BY price ASC
        """, (f'%{name}%',))
        return [dict(row) for row in self.cursor.fetchall()]
    
    def query_by_price_range(self, min_price: float, max_price: float) -> List[Dict]:
        """按价格区间查询"""
        self.cursor.execute("""
            SELECT * FROM items 
            WHERE price BETWEEN ? AND ? 
            ORDER BY price ASC
        """, (min_price, max_price))
        return [dict(row) for row in self.cursor.fetchall()]
    
    def get_cheapest(self, limit: int = 10) -> List[Dict]:
        """获取最便宜的饰品"""
        self.cursor.execute("""
            SELECT * FROM items 
            WHERE price IS NOT NULL
            ORDER BY price ASC 
            LIMIT ?
        """, (limit,))
        return [dict(row) for row in self.cursor.fetchall()]
    
    def get_stats(self) -> Dict:
        """获取统计信息"""
        self.cursor.execute("""
            SELECT 
                COUNT(*) as total,
                AVG(price) as avg_price,
                MIN(price) as min_price,
                MAX(price) as max_price,
                SUM(on_sale_count) as total_on_sale
            FROM items
        """)
        return dict(self.cursor.fetchone())
    
    def clear_all(self):
        """清空所有数据"""
        self.cursor.execute("DELETE FROM items")
        self.conn.commit()
        print("数据已清空")
    
    def close(self):
        """关闭数据库连接"""
        self.conn.close()
        print("数据库已关闭")
    
    @staticmethod
    def _to_float(value) -> Optional[float]:
        """转换为浮点数"""
        try:
            return float(value) if value else None
        except (ValueError, TypeError):
            return None

# 使用示例
def main():
    # 1. 创建数据库实例
    db = ItemDatabase("./data/data.db")
    
    # 2. 示例：插入单条数据
    sample_item = {
        "id": 110829,
        "gameId": 730,
        "gameName": "CS:GO",
        "commodityName": "MAG-7（StatTrak™） | 重新补给 (破损不堪)",
        "commodityHashName": "StatTrak™ MAG-7 | Resupply (Well-Worn)",
        "iconUrl": "https://youpin.img898.com/csgo/template/192699daca1b49c499a90dcd7b314db3.png",
        "onSaleCount": 55,
        "price": "0.59",
        "steamPrice": "1.04",
        "steamUsdPrice": "0.13",
        "typeName": "霰弹枪",
        "exterior": "破损不堪",
        "exteriorColor": "C96C69",
        "rarity": "军规级",
        "rarityColor": "7087FF",
        "quality": "StatTrak™",
        "qualityColor": "CF6A32",
        "haveLease": 1
    }
    db.insert_item(sample_item)
    
    # 3. 查询示例
    print("\n=== 按名称查询 ===")
    results = db.query_by_name("MAG-7")
    for item in results:
        print(f"{item['commodity_name']}: ¥{item['price']}")
    
    print("\n=== 价格区间查询（0.5-1元）===")
    results = db.query_by_price_range(0.5, 1.0)
    for item in results:
        print(f"{item['commodity_name']}: ¥{item['price']}")
    
    print("\n=== 最便宜的5个饰品 ===")
    cheapest = db.get_cheapest(5)
    for item in cheapest:
        print(f"{item['commodity_name']}: ¥{item['price']}")
    
    print("\n=== 统计信息 ===")
    stats = db.get_stats()
    print(f"总数: {stats['total']}")
    print(f"平均价: ¥{stats['avg_price']:.2f}")
    print(f"最低价: ¥{stats['min_price']}")
    print(f"最高价: ¥{stats['max_price']}")
    
    # 4. 关闭数据库
    db.close()

if __name__ == "__main__":
    main()