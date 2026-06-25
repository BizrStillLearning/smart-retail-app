from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import List
from recommendation import ai_engine

app = FastAPI(title="Smart Retail AI Service")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
)


class TransactionItem(BaseModel):
    id_pesanan: int
    id_produk: int


@app.get("/")
def health_check():
    return {"status": "AI Service Running Optimal"}


@app.post("/api/ai/train")
def train_association_rules(transactions: List[TransactionItem]):
    data = [{"id_pesanan": t.id_pesanan, "id_produk": t.id_produk} for t in transactions]

    success = ai_engine.train_model(data)
    if success:
        return {"message": "Model FP-Growth berhasil dilatih ulang"}
    else:
        raise HTTPException(status_code=400, detail="Data transaksi tidak cukup untuk melatih model")


@app.get("/api/ai/recommend/{id_produk}")
def get_recommendation(id_produk: int):
    rekomendasi = ai_engine.get_recommendation(id_produk)
    return {
        "id_produk_referensi": id_produk,
        "rekomendasi_id": rekomendasi
    }


