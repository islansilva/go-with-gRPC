package service

import (
	context "context"
	"io"
	"example/gRPC/internal/database"
	"example/gRPC/internal/pb"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	CategoryDB database.Category
}

func NewCategoryService(categoryDB database.Category) *CategoryService {
	return &CategoryService {
		CategoryDB: categoryDB,
	}
}


func (c *CategoryService) CreateCategory(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.CategoryResponse, error) {
	category, err := c.CategoryDB.Create(in.Name, &in.Description)

	if err != nil {
		return nil, err
	}

	return &pb.CategoryResponse {
		Category: &pb.Category {
			Id:		category.ID,
			Name: 	category.Name,
			Description:	*category.Description,
		},
	}, nil

}


func (c *CategoryService) ListCategories(ctx context.Context, blank *pb.Blank) (*pb.CategoryList, error) {
	
	categories, err := c.CategoryDB.FindAll()
	if err != nil {
		return nil, err
	}

	var categoriesResponse []*pb.Category

	for _, categories := range categories {
		categoryResponse := &pb.Category {
			Id:		categories.ID,
			Name: 	categories.Name,
			Description:	*categories.Description,
		}

		categoriesResponse = append(categoriesResponse, categoryResponse)
	}

	return &pb.CategoryList{
		Categories: categoriesResponse,
	}, nil

}


func (c *CategoryService) GetCategory(ctx context.Context, categoryGetRequest *pb.CategoryGetRequest) (*pb.Category, error) {
	categories, err := c.CategoryDB.FindByCourseID(categoryGetRequest.Id)
	if err != nil {
		return nil, err
	}

	return  &pb.Category {
		Id:		categories.ID,
		Name: 	categories.Name,
		Description:	*categories.Description,
	}, nil

}


func (c *CategoryService) CreateCategoryStream(stream pb.CategoryService_CreateCategoryStreamServer) error {
	categories := &pb.CategoryList{}

	for {
		category, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(categories)
		}

		if err != nil {
			return err
		}

		categoryResult, err := c.CategoryDB.Create(category.Name, &category.Description)
		if err != nil {
			return err
		}
		categories.Categories = append(categories.Categories, &pb.Category {
			Id:		categoryResult.ID,
			Name: 	categoryResult.Name,
			Description:	*categoryResult.Description,
		})
	}
	
}


func (c *CategoryService) CreateCategoryStreamBidirectional(stream pb.CategoryService_CreateCategoryStreamBidirectionalServer) error {
	for {
		category, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		categoryResult, err := c.CategoryDB.Create(category.Name, &category.Description)
		if err != nil {
			return err
		}
		
		err = stream.Send(&pb.Category {
			Id:		categoryResult.ID,
			Name: 	categoryResult.Name,
			Description:	*categoryResult.Description,
		})
		if err != nil {
			return err
		}

	}
}