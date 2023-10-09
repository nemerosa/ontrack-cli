package client

import "ontrack-cli/config"

func SetupPromotionLevel(
	cfg *config.Config,
	project string,
	branch string,
	promotion string,
	description string,
	autoPromotion bool,
	validations []string,
	promotions []string,
	include string,
	exclude string,
) error {

	// Data
	var data struct {
		SetupPromotionLevel struct {
			Errors []struct {
				Message string
			}
		}
		SetPromotionLevelAutoPromotionProperty struct {
			Errors []struct {
				Message string
			}
		}
	}

	// Call
	if err := GraphQLCall(cfg, `
			mutation SetupPromotionLevel(
				$project: String!,
				$branch: String!,
				$promotion: String!,
				$description: String,
				$autoPromotion: Boolean!,
				$validationStamps: [String!],
				$include: String,
				$exclude: String,
				$promotionLevels: [String!]
			) {
				setupPromotionLevel(input: {
					project: $project,
					branch: $branch,
					promotion: $promotion,
					description: $description
				}) {
					errors {
						message
					}
				}
				setPromotionLevelAutoPromotionProperty(input: {
					project: $project,
					branch: $branch,
					promotion: $promotion,
					validationStamps: $validationStamps,
					include: $include,
					exclude: $exclude,
					promotionLevels: $promotionLevels
				}) @include(if: $autoPromotion) {
					errors {
						message
					}
				}
			}
		`, map[string]interface{}{
		"project":          project,
		"branch":           branch,
		"promotion":        promotion,
		"description":      description,
		"autoPromotion":    autoPromotion,
		"validationStamps": validations,
		"promotionLevels":  promotions,
		"include":          include,
		"exclude":          exclude,
	}, &data); err != nil {
		return err
	}

	// Error checks
	if err := CheckDataErrors(data.SetupPromotionLevel.Errors); err != nil {
		return err
	}
	if err := CheckDataErrors(data.SetPromotionLevelAutoPromotionProperty.Errors); err != nil {
		return err
	}

	// OK
	return nil
}
