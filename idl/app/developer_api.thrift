namespace go app.developer_api

enum ModelClass {
    GPT             = 1
    SEED            = 2
    Claude          = 3
    MiniMax         = 4  // name: MiniMax
    Plugin          = 5
    StableDiffusion = 6
    ByteArtist      = 7
    Maas            = 9
    QianFan         = 10 // Abandoned: Qianfan (Baidu Cloud)
    Gemini          = 11 // name：Google Gemini
    Moonshot        = 12 // name: Moonshot
    GLM             = 13 // Name: Zhipu
    MaaSAutoSync    = 14 // Name: Volcano Ark
    QWen            = 15 // Name: Tongyi Qianwen
    Cohere          = 16 // name: Cohere
    Baichuan        = 17 // Name: Baichuan Intelligent
    Ernie           = 18 // Name: ERNIE Bot
    DeekSeek        = 19 // Name: Magic Square
    Llama           = 20 // name: Llama
    StepFun         = 23
    Other           = 999
}