package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/nicholeonardo13/gqlgen-todos/database"
	"github.com/nicholeonardo13/gqlgen-todos/graph/generated"
	"github.com/nicholeonardo13/gqlgen-todos/graph/middleware"
	"github.com/nicholeonardo13/gqlgen-todos/graph/model"
	"github.com/nicholeonardo13/gqlgen-todos/graph/resolver"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

const defaultPort = "2000"

func main() {

	db, err:= database.Connect()
	if err != nil{
		panic(err)
	}
	//defer db.Close()

	db.Exec("DROP TABLE game_genres ,  game_developers");
	db.Migrator().DropTable(&model.User{}, &model.Country{} , &model.Otp{} , &model.Genre{} , &model.Developer{} , &model.Publisher{} , &model.Game{})
	db.AutoMigrate(&model.User{},&model.Country{} , &model.Otp{} , &model.Genre{} , &model.Developer{} , &model.Publisher{} , &model.Game{})
	SeedUsers()
	SeedCountry()
	SeedGenre()
	SeedDeveloper()
	SeedPublisher()
	SeedGame()

	router := chi.NewRouter()

	router.Use(middleware.Middleware())

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200", "http://localhost:2000"},
		AllowCredentials: true,
		Debug:            true,

	}).Handler)


	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver.Resolver{}}))
	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Check against your desired domains here
				return r.Host == "example.org"
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})

	router.Handle("/", playground.Handler("Starwars", "/api"))
	router.Handle("/api", srv)

	err = http.ListenAndServe(":2000", router)
	if err != nil {
		panic(err)
	}
}

func SeedUsers() {
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}
	//defer db.Close()

	db.Create(&model.User{
		Username: "yoko",
		Password:    HashPassword("yoko"),
		Email:      "yoko@gmail.com",
		Money:       5000,
		CountryRegion: 1,
	})

	db.Create(&model.User{
		Username: "eren",
		Password:    HashPassword("eren"),
		Email:      "eren@gmail.com",
		Money:       7500,
		CountryRegion: 2,
	})
}

func SeedCountry(){
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}

	db.Create(&model.Country{
		CountryName: "America",
	})

	db.Create(&model.Country{
		CountryName: "China",
	})

	db.Create(&model.Country{
		CountryName: "German",
	})

	db.Create(&model.Country{
		CountryName: "Japan",
	})

	db.Create(&model.Country{
		CountryName: "Indonesia",
	})
}

func SeedGenre() {
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}

	db.Create(&model.Genre{
		GenreName: "Family",
	})

	db.Create(&model.Genre{
		GenreName: "Casual",
	})

	db.Create(&model.Genre{
		GenreName: "RPG",
	})

	db.Create(&model.Genre{
		GenreName: "FPS",
	})

	db.Create(&model.Genre{
		GenreName: "Horror",
	})

}

func SeedDeveloper() {
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}

	db.Create(&model.Developer{
		DeveloperName: "Rockstar",
	})

	db.Create(&model.Developer{
		DeveloperName: "Natsume",
	})

	db.Create(&model.Developer{
		DeveloperName: "Ubisoft",
	})

	db.Create(&model.Developer{
		DeveloperName: "EA",
	})
}

func SeedPublisher() {
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}

	db.Create(&model.Publisher{
		PublisherName: "Nintendo",
	})

	db.Create(&model.Publisher{
		PublisherName: "YG",
	})

	db.Create(&model.Publisher{
		PublisherName: "JYP",
	})

	db.Create(&model.Publisher{
		PublisherName: "SM Town",
	})
}

func SeedGame() {
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}

	db.Create(&model.Game{
		Title: "Yoko Pangestu The Game",
		Genres: []*model.Genre{
			{ID: 1},
			{ID: 3},
		},
		Tag: "Game Asik",
		Price: 1000,
		Rating: 4,
		Discount: 0,
		Description: "Yoko And The Gang",
		Developers: []*model.Developer{
			{ID: 2},
		},
		PublisherID: 2,
		ReleaseDate: time.Now(),
		Picture: "HOHO",
		Banner: "Banner HOHO",
		System: "System HOHO",
	})


	db.Create(&model.Game{
		Title: "Boba Mathca The Game",
		Genres: []*model.Genre{
			{ID: 2},
			{ID: 3},
		},
		Tag: "Game Lapar",
		Price: 2000,
		Rating: 6,
		Discount: 0,
		Description: "Si Boba",
		Developers: []*model.Developer{
			{ID: 3},
			{ID: 2},
		},
		PublisherID: 1,
		ReleaseDate: time.Now(),
		Picture: "BOBA",
		Banner: "Banner Boba",
		System: "System OP",
	})
}

func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

