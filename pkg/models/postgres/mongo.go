package postgres

import (
	"context"
	"examFortune/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type MongoModel struct {
	Pool *mongo.Database
}

//-------------------------SYLLABUS INFO _____START

func (c *MongoModel) InsertTempField(temp *models.TempFields, collectionName string) {
	_, err := c.Pool.Collection(collectionName).InsertOne(context.TODO(), temp)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *MongoModel) GetTempFields(syllabusInfoId int, collectionName string) []*models.TempFields {
	cursor, err := c.Pool.Collection(collectionName).Find(context.TODO(),
		bson.D{
			{"syllabusinfoid", syllabusInfoId},
		})
	if err != nil {
		log.Fatal(err)
	}
	var tempFields []*models.TempFields
	if err = cursor.All(context.TODO(), &tempFields); err != nil {
		log.Fatal(err)
	}

	return tempFields
}

func (c *MongoModel) DeleteTempField(temp *models.TempFields, collectionName string) {
	_, err := c.Pool.Collection(collectionName).DeleteOne(context.TODO(), bson.D{
		{"syllabusinfoid", temp.SyllabusInfoId}, {"teacherid", temp.TeacherId},
	})

	if err != nil {
		log.Fatal(err)
	}
}

func (c *MongoModel) UpdateTempField(temp *models.TempFields, collectionName string) {
	c.DeleteTempField(temp, collectionName)
	c.InsertTempField(temp, collectionName)
}

//-------------------------TABLE 1 _____START

func (c *MongoModel) InsertTempSessionTopic(topic *models.TopicWeek, collectionName string) {
	_, err := c.Pool.Collection(collectionName).InsertOne(context.TODO(), topic)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *MongoModel) GetTempSessionTopic(syllabusInfoId int, collectionName string) []*models.TopicWeek {
	cursor, err := c.Pool.Collection(collectionName).Find(context.TODO(),
		bson.D{
			{"syllabusinfoid", syllabusInfoId},
		})
	if err != nil {
		log.Fatal(err)
	}
	var tempSessionTopics []*models.TopicWeek
	if err = cursor.All(context.TODO(), &tempSessionTopics); err != nil {
		log.Fatal(err)
	}

	return tempSessionTopics
}

func (c *MongoModel) GetTempOneSessionTopic(syllabusInfoId int, weekNum int, collectionName string) *models.TopicWeek {
	cursor, err := c.Pool.Collection(collectionName).Find(context.TODO(),
		bson.D{
			{"syllabusinfoid", syllabusInfoId},
			{"weeknumber", weekNum},
		})
	if err != nil {
		log.Fatal(err)
	}

	var tempSessionTopics []*models.TopicWeek
	if err = cursor.All(context.TODO(), &tempSessionTopics); err != nil {
		log.Fatal(err)
	}

	return tempSessionTopics[0]
}

func (c *MongoModel) GetTempSessionTopicByWeek(syllabusInfoId int, weekNum int, collectionName string) []*models.TopicWeek {
	cursor, err := c.Pool.Collection(collectionName).Find(context.TODO(),
		bson.D{
			{"syllabusinfoid", syllabusInfoId},
			{"weeknumber", weekNum},
		})
	if err != nil {
		log.Fatal(err)
	}

	var tempSessionTopics []*models.TopicWeek
	if err = cursor.All(context.TODO(), &tempSessionTopics); err != nil {
		log.Fatal(err)
	}

	return tempSessionTopics
}

func (c *MongoModel) DeleteTempSessionTopic(topic *models.TopicWeek, collectionName string) {
	_, err := c.Pool.Collection(collectionName).DeleteOne(context.TODO(), bson.D{
		{"syllabusinfoid", topic.SyllabusInfoId},
		{"teacherid", topic.TeacherId},
		{"weeknumber", topic.WeekNumber},
	})

	if err != nil {
		log.Fatal(err)
	}
}

func (c *MongoModel) UpdateTempSessionTopic(topic *models.TopicWeek, collectionName string) {
	c.DeleteTempSessionTopic(topic, collectionName)
	c.InsertTempSessionTopic(topic, collectionName)
}

//-------------------------TABLE 2 _____START

func (c *MongoModel) InsertTempStudentTopic(topic *models.StudentTopicWeek, collectionName string) {
	_, err := c.Pool.Collection(collectionName).InsertOne(context.TODO(), topic)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *MongoModel) GetTempStudentTopic(syllabusInfoId int, collectionName string) []*models.StudentTopicWeek {
	cursor, err := c.Pool.Collection(collectionName).Find(context.TODO(),
		bson.D{
			{"syllabusinfoid", syllabusInfoId},
		})
	if err != nil {
		log.Fatal(err)
	}
	var tempStudentTopics []*models.StudentTopicWeek
	if err = cursor.All(context.TODO(), &tempStudentTopics); err != nil {
		log.Fatal(err)
	}

	return tempStudentTopics
}

func (c *MongoModel) GetTempOneStudentTopic(syllabusInfoId int, weekNum int, collectionName string) *models.StudentTopicWeek {
	cursor, err := c.Pool.Collection(collectionName).Find(context.TODO(),
		bson.D{
			{"syllabusinfoid", syllabusInfoId},
			{"weeknumber", weekNum},
		})
	if err != nil {
		log.Fatal(err)
	}
	var tempStudentTopics []*models.StudentTopicWeek
	if err = cursor.All(context.TODO(), &tempStudentTopics); err != nil {
		log.Fatal(err)
	}

	return tempStudentTopics[0]
}

func (c *MongoModel) GetTempStudentTopicByWeek(syllabusInfoId int, weekNum int, collectionName string) []*models.StudentTopicWeek {
	cursor, err := c.Pool.Collection(collectionName).Find(context.TODO(),
		bson.D{
			{"syllabusinfoid", syllabusInfoId},
			{"weeknumber", weekNum},
		})
	if err != nil {
		log.Fatal(err)
	}

	var tempSessionTopics []*models.StudentTopicWeek
	if err = cursor.All(context.TODO(), &tempSessionTopics); err != nil {
		log.Fatal(err)
	}

	return tempSessionTopics
}

func (c *MongoModel) DeleteTempStudentTopic(topic *models.StudentTopicWeek, collectionName string) {
	_, err := c.Pool.Collection(collectionName).DeleteOne(context.TODO(), bson.D{
		{"syllabusinfoid", topic.SyllabusInfoId},
		{"teacherid", topic.TeacherId},
		{"weeknumber", topic.WeekNumber},
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (c *MongoModel) UpdateTempStudentTopic(topic *models.StudentTopicWeek, collectionName string) {
	c.DeleteTempStudentTopic(topic, collectionName)
	c.InsertTempStudentTopic(topic, collectionName)
}
