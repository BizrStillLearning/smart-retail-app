import pandas as pd
from mlxtend.frequent_patterns import fpgrowth, association_rules


class FPGrowthEngine:
    def __init__(self):
        self.rules = None

    def train_model(self, transactions_data):
        # [{"id_pesanan": 1, "id_produk": 2}, {"id_pesanan": 1, "id_produk": 5}, ...]

        if not transactions_data:
            return False

        df = pd.DataFrame(transactions_data)

        basket = (df.groupby(['id_pesanan', 'id_produk'])['id_produk']
                  .count().unstack().reset_index().fillna(0)
                  .set_index('id_pesanan'))

        basket_encoded = basket.applymap(lambda x: 1 if x > 0 else 0)

        frequent_itemsets = fpgrowth(basket_encoded, min_support=0.05, use_colnames=True)

        if not frequent_itemsets.empty:
            self.rules = association_rules(frequent_itemsets, metric="lift", min_threshold=1.0)
            return True
        return False

    def get_recommendation(self, item_id):
        if self.rules is None or self.rules.empty:
            return []

        recommendations = self.rules[self.rules['antecedents'] == frozenset({item_id})]

        recommendations = recommendations.sort_values(by='confidence', ascending=False)

        res = [list(x)[0] for x in recommendations['consequents'].values]
        return res


ai_engine = FPGrowthEngine()