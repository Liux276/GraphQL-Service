package GraphQL_Service

import (
	// "encoding/json"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"strings"

	// "github.com/Go-GraphQL-Group/GraphQL-Service/server/service"
)

type Resolver struct{}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
func (r *queryResolver) People(ctx context.Context, id string) (*People, error) {
	err, people := GetPeopleByID(id)
	checkErr(err)
	return people, err
}
func (r *queryResolver) Film(ctx context.Context, id string) (*Film, error) {
	err, film := GetFilmByID(id)
	checkErr(err)
	return film, err
}
func (r *queryResolver) Starship(ctx context.Context, id string) (*Starship, error) {
	err, starship := GetStarshipByID(id)
	checkErr(err)
	return starship, err
}
func (r *queryResolver) Vehicle(ctx context.Context, id string) (*Vehicle, error) {
	err, vehicle := GetVehicleByID(id)
	checkErr(err)
	return vehicle, err
}
func (r *queryResolver) Specie(ctx context.Context, id string) (*Specie, error) {
	err, specie := GetSpeciesByID(id)
	checkErr(err)
	return specie, err
}
func (r *queryResolver) Planet(ctx context.Context, id string) (*Planet, error) {
	err, planet := GetPlanetByID(id)
	checkErr(err)
	return planet, err
}

func encodeCursor(i uint) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("cursor%d", i)))
}

func (r *queryResolver) Peoples(ctx context.Context, first *int, after *string) (PeopleConnection, error) {
	
	from := 1
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return PeopleConnection{}, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return PeopleConnection{}, err
		}
		from = i + 1
	}
	count := first
	// 获取edges
	edges := []PeopleEdge{}

	var i uint

	for i = uint(from); i < uint(to); i++ {
		_, people := GetPeopleByID(strconv.Itoa(int(i)))
		if people.ID == "" {
			break
		}
		edges = append(edges, PeopleEdge{
			Node:   people,
			Cursor: encodeCursor(i),
		})
	}

	// 获取pageInfo
	pageInfo := PageInfo{
		StartCursor: encodeCursor(uint(from)),
	}
	if from == 1 || len(edges) == 0 {
		pageInfo.HasPreviousPage = false
	} else {
		pageInfo.HasPreviousPage = true
	}
	if i < uint(to) || len(edges) == 0 {
		pageInfo.HasNextPage = false
	} else if i == uint(to) {
		_, people := GetPeopleByID(strconv.Itoa(int(i)))
		if people.ID == "" {
			pageInfo.HasNextPage = false
		} else {
			pageInfo.HasNextPage = true
		}
	}
	if len(edges) == 0 {
		pageInfo.EndCursor = encodeCursor(i)
	} else {
		pageInfo.EndCursor = encodeCursor(i - 1)
	}
	return PeopleConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: int(i) - from,
	}, nil
}
func (r *queryResolver) Films(ctx context.Context, first *int, after *string) (FilmConnection, error) {
	panic("not implemented")
}
func (r *queryResolver) Starships(ctx context.Context, first *int, after *string) (StarshipConnection, error) {
	panic("not implemented")
}
func (r *queryResolver) Vehicles(ctx context.Context, first *int, after *string) (VehicleConnection, error) {
	panic("not implemented")
}
func (r *queryResolver) Species(ctx context.Context, first *int, after *string) (SpecieConnection, error) {
	panic("not implemented")
}
func (r *queryResolver) Planets(ctx context.Context, first *int, after *string) (PlanetConnection, error) {
	panic("not implemented")
}

func (r *queryResolver) PeopleSearch(ctx context.Context, search string, first *int, after *string) (*PeopleConnection, error) {
	fmt.Println(ctx)
	// token := &service.Token{}
	// tokenJson, _ := ctx.Value(service.Issuer).(string)
	// json.Unmarshal([]byte(tokenJson), token)
	// service.ParseToken(token.SW_TOKEN, []byte(service.SecretKey))

	from := 1
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return &PeopleConnection{}, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return &PeopleConnection{}, err
		}
		from = i + 1
	}
	hasPreviousPage := false
	for i := 1; i <= from; i++ {
		_, people := GetPeopleByID(strconv.Itoa(from))
		if people.ID == "" {
			break
		}
		if people.Name == search {
			hasPreviousPage = true
			break
		}
	}

	// 获取edges
	edges := []PeopleEdge{}
	for {
		_, people := GetPeopleByID(strconv.Itoa(from))
		if people.ID == "" {
			break
		}
		if people.Name == search {
			edges = append(edges, PeopleEdge{
				Node:   people,
				Cursor: encodeCursor(uint(from)),
			})
		}
		from++
		if len(edges) == *first {
			break
		}
	}

	// 获取pageInfo
	startId, _ := strconv.Atoi(edges[0].Node.ID)
	endId, _ := strconv.Atoi(edges[0].Node.ID)
	pageInfo := PageInfo{
		StartCursor:     encodeCursor(uint(startId)),
		EndCursor:       encodeCursor(uint(endId)),
		HasPreviousPage: hasPreviousPage,
	}
	hasNextPage := false
	if len(edges) == *first {
		for i := endId + 1; ; i++ {
			_, people := GetPeopleByID(strconv.Itoa(i))
			if people.ID == "" {
				break
			}
			if people.Name == search {
				hasNextPage = true
				break
			}
		}
	}
	pageInfo.HasNextPage = hasNextPage
	return &PeopleConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: len(edges),
	}, nil
}
func (r *queryResolver) FilmsSearch(ctx context.Context, search string, first *int, after *string) (*FilmConnection, error) {
	panic("not implemented")
}
func (r *queryResolver) StarshipsSearch(ctx context.Context, search string, first *int, after *string) (*StarshipConnection, error) {
	panic("not implemented")
}
func (r *queryResolver) VehiclesSearch(ctx context.Context, search string, first *int, after *string) (*VehicleConnection, error) {
	panic("not implemented")
}
func (r *queryResolver) SpeciesSearch(ctx context.Context, search string, first *int, after *string) (*SpecieConnection, error) {
	panic("not implemented")
}
func (r *queryResolver) PlanetsSearch(ctx context.Context, search string, first *int, after *string) (*PlanetConnection, error) {
	panic("not implemented")
}
